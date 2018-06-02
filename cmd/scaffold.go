package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:     "scaffold [name] atribute:type!",
	Aliases: []string{"scaf", "s", "Scaffold"},
	Short:   "Creates attribute new object",
	Long: `Scaffold (goals scaffold) will create attribute new object
	and it's structure, nammed: Model. Schema and Resolvers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			er("Wrong arguments you should use a minimum of 2 arguments")
		}
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}

		project := RecreateProjectFromGoals(wd)

		createFiles(args[0], args[1:], project)
	},
}

func createFiles(name string, args []string, project Project) {
	model, schema, methods := getTemplates(args)
	resolver := fmt.Sprintf("%s%sResolver", strings.ToLower(string(name[0])), name[1:])
	abbreviation := toAbbreviation(name)
	name = strings.Title(name)

	data := map[string]string{"model": model, "schema": schema, "Name": name, "abbreviation": abbreviation, "resolver": resolver}

	methods = replaceTemplate(methods, data)

	data["methods"] = methods

	modelScript := replaceTemplate(Templates["fullmodel"], data)
	schemaScript := executeTemplate(Templates["fullschema"], data)
	resolverScript := replaceTemplate(Templates["fullresolver"], data)

	writeStringToFile(filepath.Join("app/model", fmt.Sprintf("%s.go", strings.ToLower(name))), modelScript)
	writeStringToFile(filepath.Join("app/schema", fmt.Sprintf("%ssch.go", strings.ToLower(name))), schemaScript)
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

func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func toAbbreviation(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	out = append(out, unicode.ToLower(runes[0]))
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, unicode.ToLower(runes[i]))
		}
	}

	return string(out)
}
