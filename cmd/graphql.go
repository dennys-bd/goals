package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dennys-bd/goals/core"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	"github.com/spf13/cobra"
)

var resolver string
var json, details, noGqllModel, noGqlSchema, noGqlResolver bool

var gqlCmd = &cobra.Command{
	Use:     "graphql Name 'atribute:type!'",
	Aliases: []string{"gql", "g"},
	Short:   "Creates graphql model's structure",
	Long: `Graphql (goals scaffold graphql or symply goals s g)
based on a model's description will create
it's structure, nammed: Model, Schema and Resolver.`,
	Run: func(cmd *cobra.Command, args []string) {

		println("len of args: ", len(args))
		if len(args) < 1 || (len(args) < 2 && !details) {
			errs.Ex("Wrong arguments you should insert at least the name of your model, and it's structure if flag details is not setted")
		}

		project := recreateProjectFromGoals()
		gqlTemplates()
		createFiles(args[0], args[1:], project)

	},
}

func createFiles(name string, args []string, project core.Project) {
	model, schema, resolver, modelMethods := getTemplates(args)
	name = strings.Title(name)
	if !noGqllModel {
		writeModel(name, modelMethods, model, project)
	}
	if !noGqlSchema {
		writeSchema(name, schema, project)
	}
	if !noGqlResolver {
		writeResolver(name, resolver, project)
	}
}

func init() {
	gqlCmd.Flags().StringVarP(&resolver, "resolver", "r", "", "Name to the resolver variable of your model")
	gqlCmd.Flags().BoolVar(&json, "json", false, "Use it if you want to generate the json attributes of your model")
	gqlCmd.Flags().BoolVar(&noGqllModel, "no-model", false, "Use this flag if you want to graphql command to don't create the model")
	gqlCmd.Flags().BoolVar(&noGqlSchema, "no-schema", false, "Use this flag if you want to graphql command to don't create the schema")
	gqlCmd.Flags().BoolVar(&noGqlResolver, "no-resolver", false, "Use this flag if you want to graphql command to don't create the resolver")
	gqlCmd.Flags().BoolVarP(&details, "details", "d", false, "Details will give the opportunity create your model line by line")
}

func writeModel(name, methods, template string, project core.Project) {
	importpath := ""
	countImports := 0
	if strings.Index(template, "graphql.ID") > -1 {
		importpath += "\"github.com/dennys-bd/goals/graphql\"\n"
		countImports++
	}
	if strings.Index(template, "scalar.") > -1 {
		importpath += fmt.Sprintf("	\"%v/app/scalar\"\n", project.ImportPath)
		countImports++
	}
	if strings.Index(template, "time.Time") > -1 {
		importpath = "	\"time\"\n\n"
		countImports++
	}

	data := make(map[string]interface{})

	if countImports > 1 {
		importpath = fmt.Sprintf("import (\n%v)\n\n", importpath)
		data = map[string]interface{}{"model": template, "Name": name, "importpath": importpath}
	} else if countImports == 1 {
		importpath = "import " + importpath
		data = map[string]interface{}{"model": template, "Name": name, "importpath": importpath}
	} else {
		data = map[string]interface{}{"model": template, "Name": name}
	}

	modelScript := executeTemplate(templates["scafmodel"], data)
	modelScript += methods
	modelScript = strings.Replace(modelScript, "&#34;", "\"", -1)
	modelScript = strings.Replace(modelScript, "&amp;", "&", -1)
	modelScript = strings.Replace(modelScript, "model_name", name, -1)
	modelScript = strings.Replace(modelScript, "&#39;", "'", -1)

	writeStringToFile(filepath.Join("app/model", fmt.Sprintf("%s.go", strings.ToLower(name))), modelScript)
}

func writeSchema(name string, template string, project core.Project) {
	data := map[string]string{"schema": template, "Name": name, "name": strings.ToLower(string(name[0])) + name[1:]}
	schemaScript := executeTemplate(templates["scafschema"], data)
	writeStringToFile(filepath.Join("app/schema", fmt.Sprintf("%s_schema.go", strings.ToLower(name))), schemaScript)
}

func writeResolver(name string, template string, project core.Project) {
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

	writeStringToFile(filepath.Join("app/resolver", fmt.Sprintf("%s_resolver.go", strings.ToLower(name))), resolverScript)
}

func getTemplates(args []string) (model, schema, resolver, modelMethods string) {

	var mB, sB, rB, mmB bytes.Buffer
	if details {
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Printf("Create your model (%s) based on a GraphQL's structure:\n\n", name)
		for {
			fmt.Printf("\rEnter your model's next attribute")
			// text, _ := reader.ReadString('\n')

		}
	}
	if len(args) == 1 {
		args = strings.Split(args[0], " ")
	}

	for _, arg := range args {
		attr, tyName, model, mandatory, list, manL := getLineAttributes(arg)
		mB.WriteString(getModelLine(attr, tyName, model, mandatory, list, manL))
		sB.WriteString(getSchemaLine(attr, tyName, mandatory, list, manL))
		rB.WriteString(getResolverLine(attr, tyName, model, mandatory, list, manL))
		mmB.WriteString(getModelMethods(attr, tyName, model, mandatory, list, manL))
	}
	return mB.String(), sB.String(), rB.String(), mmB.String()
}

func getLineAttributes(argument string) (attribute, typeName string, isModel, isMandatory, isList, isMandatoryInList bool) {

	arguments := strings.Split(argument, ":")
	if len(arguments) == 2 {
		attribute, typeName, isModel = arguments[0], arguments[1], false
	} else if len(arguments) == 3 && arguments[1] == "type" {
		attribute, typeName, isModel = arguments[0], arguments[2], true
	} else {
		errs.Ex(fmt.Sprintf("Error: Bad Syntax in %s", argument))
	}

	isMandatory = strings.HasSuffix(typeName, "!")
	isList = strings.HasPrefix(typeName, "[")
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
			errs.Ex(fmt.Sprintf("Bad Syntax: %s should close list with ]", typeName))
		}
	} else if strings.HasSuffix(typeName, "]") {
		errs.Ex(fmt.Sprintf("Bad Syntax: %s should start list with [ before close", typeName))
	}

	return attribute, typeName, isModel, isMandatory, isList, isMandatoryInList
}

func getModelLine(attribute, typeName string, isModel, isMandatory, isList, isMandatoryInList bool) string {
	if strings.EqualFold(attribute, "id") {
		attribute = strings.ToUpper(attribute)
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

	if json {
		if isModel {
			return fmt.Sprintf("	%s %s `json:\"-\"`\n", strings.Title(attribute), strings.Title(typeName))
		} else if typeName == "*bool" || typeName == "bool" {
			return fmt.Sprintf("	%s %s `json:\"%s\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
		} else {
			return fmt.Sprintf("	%s %s `json:\"%s,omitempty\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
		}
	}
	if isModel {
		return fmt.Sprintf("	%s %s", strings.Title(attribute), strings.Title(typeName))
	}
	return fmt.Sprintf("	%s %s", strings.Title(attribute), typeName)
}

func getSchemaLine(attribute, typeName string, isMandatory, isList, isMandatoryInList bool) string {
	typeReturn := typeName
	switch typeName {
	case "boolean", "Bool", "bool":
		typeReturn = "Boolean"
	case "id", "Id":
		typeReturn = "ID"
	case "time", "Time":
		typeReturn = "String"
		typeName = "time"
	default:
		typeReturn = strings.Title(typeReturn)
	}

	if isList {
		if isMandatoryInList {
			typeReturn = "[" + typeReturn + "!]"
		} else {
			typeReturn = "[" + typeReturn + "]"
		}
	}

	if isMandatory {
		typeReturn += "!"
	}

	if typeName != "time" {
		return fmt.Sprintf("	%s: %s\n", attribute, typeReturn)
	}
	return fmt.Sprintf("	%s(format: String): %s\n", attribute, typeReturn)
}

func getResolverLine(attribute, typeName string, isModel, isMandatory, isList, isMandatoryInList bool) string {
	if !isModel {
		var typeReturn string
		switch typeName {
		case "String", "string":
			typeReturn = "string"
		case "Int", "int":
			typeReturn = "int32"
		case "Float", "float":
			typeReturn = "float64"
		case "Boolean", "boolean", "Bool", "bool":
			typeReturn = "bool"
		case "ID", "id":
			typeReturn = "graphql.ID"
		case "time", "Time":
			typeReturn = "string"
			typeName = "time"
		default:
			typeReturn = fmt.Sprintf("scalar.%s", strings.Title(typeName))
		}

		if isList {
			if isMandatoryInList {
				typeReturn = "[]" + typeReturn
			} else {
				typeReturn = "[]*" + typeReturn
			}
		}

		if !isMandatory {
			typeReturn = "*" + typeReturn
		}

		if typeName != "time" {
			return fmt.Sprintf(`func (r *{{.resolver}}) %s() %s {
	return r.{{.abbreviation}}.%s
}
`, strings.Title(attribute), typeReturn, strings.Title(attribute))
		}

		return fmt.Sprintf(`func (r *{{.resolver}}) %s(args struct{ Format *string }) %s {
	return r.{{.abbreviation}}.Get%s(args.Format)
}
`, strings.Title(attribute), typeReturn, strings.Title(attribute))
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

func getModelMethods(attribute, typeName string, isModel, isMandatory, isList, isMandatoryInList bool) string {
	if !isModel {
		switch typeName {
		case "time", "Time":
			data := map[string]interface{}{"attribute": attribute, "abbreviation": toAbbreviation(attribute), "Attribute": strings.Title(attribute), "type": "%s"}
			if !isMandatory {
				data["notMandatory"] = "*"
			}
			if isList {
				data["list"] = "[]"
				if !isMandatoryInList {
					data["notInList"] = "*"
				}
			}
			return executeTemplate(templates["getDate"], data)
		}
	}

	return ""
}

func separateLine(arg string) {
	// r := regexp.MustCompile(`(?P<attribute>[a-z][a-zA-Z0-9]*)(?P<params>\([[a-z][a-zA-Z0-9]*\: ?[a-zA-Z0-9]+!?\)]*)?\: ?(?P<result>([a-z][a-zA-Z0-9]*!?)|(\[[a-z][a-zA-Z0-9]*!?\]!?))`)
}
