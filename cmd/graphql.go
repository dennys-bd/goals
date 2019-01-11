package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var resolver string

var gqlCmd = &cobra.Command{
	Use:     "graphql Name atribute:type!",
	Aliases: []string{"gql", "g"},
	Short:   "Creates graphql model's structure",
	Long: `Graphql (goals scaffold graphql or symply goals s g)
based on a model's description will create
it's structure, nammed: Model, Schema and Resolver.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			er("Wrong arguments you should use a minimum of 2 arguments")
		}

		project := recreateProjectFromGoals()

		createFiles(args[0], args[1:], project)
	},
}

func createFiles(name string, args []string, project Project) {
	model, schema, resolver := getTemplates(args)
	name = strings.Title(name)
	gqlTemplates()
	writeModel(name, model, project)
	writeSchema(name, schema, project)
	writeResolver(name, resolver, project)

}

func init() {
	gqlCmd.Flags().StringVarP(&resolver, "resolver", "r", "", "Name to the resolver variable of your model")
}

func writeModel(name string, template string, project Project) {
	hasGraphql := strings.Index(template, "graphql.ID") > -1
	hasScalar := strings.Index(template, "scalar.") > -1
	hasTime := strings.Index(template, "time.Time") > -1

	importpath := ""
	if hasTime {
		importpath = "	\"time\"\n\n"
	}
	if hasScalar {
		importpath += fmt.Sprintf("	\"%v/app/scalar\"\n", project.ImportPath)
	}
	if hasGraphql {
		importpath += "\"github.com/dennys-bd/goals/graphql\"\n"
	}

	if (hasTime && hasScalar) || (hasTime && hasGraphql) || (hasScalar && hasGraphql) {
		importpath = fmt.Sprintf("import (\n%v)\n\n", importpath)
	} else {
		importpath = "import " + importpath + "\n"
	}

	data := map[string]string{"model": template, "Name": name, "importpath": importpath}

	modelScript := replaceTemplate(templates["scafmodel"], data)

	writeStringToFile(filepath.Join("app/model", fmt.Sprintf("%s.go", strings.ToLower(name))), modelScript)
}

func writeSchema(name string, template string, project Project) {
	data := map[string]string{"schema": template, "Name": name}
	schemaScript := executeTemplate(templates["scafschema"], data)
	writeStringToFile(filepath.Join("app/schema", fmt.Sprintf("%ssch.go", strings.ToLower(name))), schemaScript)
}

func writeResolver(name string, template string, project Project) {
	if resolver == "" {
		resolver = fmt.Sprintf("%s%sResolver", strings.ToLower(string(name[0])), name[1:])
	}
	abbreviation := toAbbreviation(name)

	hasGraphql := strings.Index(template, "graphql.ID") > -1
	hasScalar := strings.Index(template, "scalar.") > -1
	hasTime := strings.Index(template, "time.Time") > -1

	importpath := ""
	if hasTime {
		importpath = "	\"time\"\n\n"
	}
	importpath = fmt.Sprintf("%s	\"%s/app/model\"\n", importpath, project.ImportPath)
	if hasScalar {
		importpath += fmt.Sprintf("	\"%v/app/scalar\"\n", project.ImportPath)
	}
	if hasGraphql {
		importpath += "	graphql \"github.com/graph-gophers/graphql-go\"\n"
	}

	if hasTime || hasScalar || hasGraphql {
		importpath = fmt.Sprintf("import (\n%v)\n\n", importpath)
	} else {
		importpath = "import " + importpath + "\n"
	}

	data := map[string]string{"Name": name, "abbreviation": abbreviation, "resolver": resolver, "importpath": importpath}

	template = replaceTemplate(template, data)

	data["methods"] = template

	resolverScript := replaceTemplate(templates["scafresolver"], data)

	writeStringToFile(filepath.Join("app/resolver", fmt.Sprintf("%ss.go", strings.ToLower(name))), resolverScript)
}

func getTemplates(args []string) (model string, schema string, resolver string) {

	var mB, sB, rB bytes.Buffer

	for _, attribute := range args {
		arguments := strings.Split(attribute, ":")
		if len(arguments) == 2 {
			attr, tyName := arguments[0], arguments[1]
			mB.WriteString(GetModelLine(attr, tyName, false))
			sB.WriteString(GetSchemaLine(attr, tyName))
			rB.WriteString(GetResolverLine(attr, tyName, false))

		} else if len(arguments) == 3 && arguments[1] == "type" {
			attr, tyName := arguments[0], arguments[2]
			mB.WriteString(GetModelLine(attr, tyName, true))
			sB.WriteString(GetSchemaLine(attr, tyName))
			rB.WriteString(GetResolverLine(attr, tyName, true))
		} else {
			er(fmt.Sprintf("Bad Syntax in %s", attribute))
		}
	}

	return mB.String(), sB.String(), rB.String()
}

// GetModelLine returns a line for model struct based on
// arguments comming from goals scaffold
func GetModelLine(attribute string, typeName string, isModel bool) string {
	if strings.EqualFold(attribute, "id") {
		attribute = strings.ToUpper(attribute)
	}
	var isMandatoryInList bool
	isMandatory := strings.HasSuffix(typeName, "!")
	isList := strings.HasPrefix(typeName, "[")
	if isMandatory {
		typeName = typeName[:len(typeName)-1]
	}

	if isList {
		if strings.HasSuffix(typeName, "]") {
			typeName = typeName[1 : len(typeName)-1]
			isMandatoryInList = strings.HasSuffix(typeName, "!")
			if isMandatoryInList {
				typeName = typeName[:len(typeName)-1]
			}
		} else {
			er(fmt.Sprintf("Bad Syntax: %s should close list with ]", typeName))
		}
	} else if strings.HasSuffix(typeName, "]") {
		er(fmt.Sprintf("Bad Syntax: %s should start list with [ before close", typeName))
	}

	if !isModel {
		switch typeName {
		case "String", "string":
			typeName = "string"
		case "Int", "int":
			typeName = "int32"
		case "Float", "float":
			typeName = "float64"
		case "Boolean", "boolean", "Bool", "bool":
			typeName = "bool"
		case "ID", "id":
			typeName = "graphql.ID"
		case "time", "Time":
			typeName = "time.Time"
		default:
			typeName = fmt.Sprintf("scalar.%s", strings.Title(typeName))
		}
	}

	if isList {
		if isMandatoryInList {
			typeName = "[]" + typeName
		} else {
			typeName = "[]*" + typeName
		}
	}

	if !isMandatory {
		typeName = "*" + typeName
	}

	if isModel {
		return fmt.Sprintf("	%s %s `json:\"-\"`\n", strings.Title(attribute), strings.Title(typeName))
	} else if typeName == "*bool" || typeName == "bool" {
		return fmt.Sprintf("	%s %s `json:\"%s\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	} else {
		return fmt.Sprintf("	%s %s `json:\"%s,omitempty\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	}
}

// GetSchemaLine returns a line for model schema based on
// arguments comming from goals scaffold
func GetSchemaLine(attribute string, typeName string) string {
	var isMandatoryInList bool
	isMandatory := strings.HasSuffix(typeName, "!")
	isList := strings.HasPrefix(typeName, "[")
	if isMandatory {
		typeName = typeName[:len(typeName)-1]
	}

	if isList {
		if strings.HasSuffix(typeName, "]") {
			typeName = typeName[1 : len(typeName)-1]
			isMandatoryInList = strings.HasSuffix(typeName, "!")
			if isMandatoryInList {
				typeName = typeName[:len(typeName)-1]
			}
		}
	}

	switch typeName {
	case "boolean", "Bool", "bool":
		typeName = "Boolean"
	case "id":
		typeName = "ID"
	default:
		typeName = strings.Title(typeName)
	}

	if isList {
		if isMandatoryInList {
			typeName = "[" + typeName + "!]"
		} else {
			typeName = "[" + typeName + "]"
		}
	}

	if isMandatory {
		typeName += "!"
	}

	return fmt.Sprintf("	%s: %s\n", attribute, typeName)
}

// GetResolverLine returns a line for model resolver based on
// arguments comming from goals scaffold
func GetResolverLine(attribute string, typeName string, isModel bool) string {
	if strings.EqualFold(attribute, "id") {
		attribute = strings.ToUpper(attribute)
	}
	var isMandatoryInList bool
	isMandatory := strings.HasSuffix(typeName, "!")
	isList := strings.HasPrefix(typeName, "[")
	if isMandatory {
		typeName = typeName[:len(typeName)-1]
	}

	if isList {
		if strings.HasSuffix(typeName, "]") {
			typeName = typeName[1 : len(typeName)-1]
			isMandatoryInList = strings.HasSuffix(typeName, "!")
			if isMandatoryInList {
				typeName = typeName[:len(typeName)-1]
			}
		}
	}

	if !isModel {
		switch typeName {
		case "String", "string":
			typeName = "string"
		case "Int", "int":
			typeName = "int32"
		case "Float", "float":
			typeName = "float64"
		case "Boolean", "boolean", "Bool", "bool":
			typeName = "bool"
		case "ID", "id":
			typeName = "graphql.ID"
		case "time", "Time":
			typeName = "time.Time"
		default:
			typeName = fmt.Sprintf("scalar.%s", strings.Title(typeName))
		}

		if isList {
			if isMandatoryInList {
				typeName = "[]" + typeName
			} else {
				typeName = "[]*" + typeName
			}
		}

		if !isMandatory {
			typeName = "*" + typeName
		}

		return fmt.Sprintf(`func (r *{{.resolver}}) %s() %s {
	return r.{{.abbreviation}}.%s
}
`, strings.Title(attribute), typeName, strings.Title(attribute))
	}

	typeName = fmt.Sprintf("%s%sResolver", strings.ToLower(string(typeName[0])), typeName[1:])

	if isList {
		pointer := "*"
		address := "&"
		insideAddress := ""
		check := `if r.{{.abbreviation}}.{{.attribute}} == nil {
		return nil
	}
	`
		bal := `func (r *{{.resolver}}) {{.attribute}}() {{.pointer}}[]*{{.typeName}} {
	{{.check}}slice := {{.pointer}}r.{{.abbreviation}}.{{.attribute}}

	l := make([]*{{.typeName}}, len(slice))
	for i := range l {
		l[i] = &{{.typeName}}{{{.insideAddress}}slice[i]}
	}

	return {{.address}}l
}
`
		if isMandatory {
			pointer = ""
			address = ""
			check = ""
		}

		if isMandatoryInList {
			insideAddress = "&"
		}

		bal = strings.Replace(bal, "{{.check}}", check, -1)
		bal = strings.Replace(bal, "{{.insideAddress}}", insideAddress, -1)
		bal = strings.Replace(bal, "{{.attribute}}", strings.Title(attribute), -1)
		bal = strings.Replace(bal, "{{.typeName}}", typeName, -1)
		bal = strings.Replace(bal, "{{.pointer}}", pointer, -1)
		bal = strings.Replace(bal, "{{.address}}", address, -1)

		return bal

	}

	address := ""
	if isMandatory {
		address = "&"
	}

	return fmt.Sprintf(`func (r *{{.resolver}}) %s() *%s {
	return &%s{%sr.{{.abbreviation}}.%s}
}
`, strings.Title(attribute), typeName, typeName, address, strings.Title(attribute))
}
