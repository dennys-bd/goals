package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dennys-bd/goals/core"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	"github.com/spf13/cobra"
)

var resolver string
var json, gqlVerbose, noGqllModel, noGqlSchema, noGqlResolver bool

type attribute struct {
	name              string
	typeName          string
	isModel           bool
	isMandatory       bool
	isList            bool
	isMandatoryInList bool
	params            []attribute
}

// attribute, typeName string, isModel, isMandatory, isList, isMandatoryInList bool

var gqlCmd = &cobra.Command{
	Use:     "graphql Name 'atribute:type!'",
	Aliases: []string{"gql", "g"},
	Short:   "Creates graphql model's structure",
	Long: `Graphql (goals scaffold graphql or symply goals s g)
based on a model's description will create
it's structure, nammed: Model, Schema and Resolver.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || (len(args) < 2 && !gqlVerbose) {
			errs.Ex("Wrong arguments you should insert at least the name of your model, and it's structure if flag verbose is not setted")
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
	gqlCmd.Flags().BoolVarP(&gqlVerbose, "verbose", "v", false, "Verbose will give the opportunity create your model line by line")
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
		importpath = "	\"time\"\n"
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
	importpath = fmt.Sprintf("%s \"%s/app/model\"\n", importpath, project.ImportPath)
	if hasScalar {
		importpath += fmt.Sprintf("	\"%v/app/scalar\"\n", project.ImportPath)
	}
	if hasGraphql {
		importpath += "	\"github.com/dennys-bd/goals/graphql\"\n"
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
	if gqlVerbose {
		r := bufio.NewReader(os.Stdin)
		fmt.Printf("Create your model (%s) based on a GraphQL's structure:\n\n", "bla")
		for {
			fmt.Printf("Enter your model's next attribute, or just press enter to finish creation. ")
			text, _ := r.ReadString('\n')
			if strings.TrimSpace(text) == "" {
				break
			}
			separateLine(text)
			attr := getLineAttributes(text)
			ml := getModelLine(attr)
			mml := getModelMethods(attr)
			sl := getSchemaLine(attr)
			rl := getResolverLine(attr)

			fmt.Printf("ModelLine: %s", ml)
			if mml != "" {
				fmt.Printf("ModelMethods: \n%s\n", strings.Trim(mml, "\n"))
			}
			fmt.Printf("SchemaLine: %s\n", strings.Trim(sl, "\n"))
			fmt.Printf("ResolverLine: \n%s\n", rl)

			print("Is it correct? (Y/n) ")
			ru, _, _ := r.ReadRune()
			if string(ru) == "y" || string(ru) == "Y" {
				mB.WriteString(ml)
				mmB.WriteString(mml)
				sB.WriteString(sl)
				rB.WriteString(rl)
			}
			println("")
			r.Reset(os.Stdin)
		}
	} else {
		if len(args) == 1 {
			args = strings.Split(args[0], " ")
		}

		for _, arg := range args {
			attr := getLineAttributes(arg)
			mB.WriteString(getModelLine(attr))
			sB.WriteString(getSchemaLine(attr))
			rB.WriteString(getResolverLine(attr))
			mmB.WriteString(getModelMethods(attr))
		}
	}
	return mB.String(), sB.String(), rB.String(), mmB.String()
}

func getLineAttributes(argument string) attribute {
	var name, typeName string
	var isModel, isMandatory, isList, isMandatoryInList bool
	var params []attribute

	opening := strings.Index(argument, "(")
	if opening > -1 {
		closing := strings.Index(argument, ")")
		p := argument[opening+1 : closing]
		argument = argument[:opening] + argument[closing+1:]
		pl := strings.Split(p, ",")
		for _, l := range pl {
			a := getLineAttributes(l)
			params = append(params, a)
		}
	}

	arguments := strings.Split(argument, ":")
	if len(arguments) == 2 {
		name, typeName, isModel = strings.TrimSpace(arguments[0]), strings.TrimSpace(arguments[1]), false
	} else if len(arguments) == 3 && arguments[1] == "type" {
		name, typeName, isModel = strings.TrimSpace(arguments[0]), strings.TrimSpace(arguments[2]), true
	} else {
		errs.Ex(fmt.Sprintf("Bad Syntax in %s", argument))
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
			errs.Ex(fmt.Sprintf("Bad Syntax, %s should close list with ]", typeName))
		}
	} else if strings.HasSuffix(typeName, "]") {
		errs.Ex(fmt.Sprintf("Bad Syntax, %s should start list with [ before close", typeName))
	}

	a := attribute{
		name:              name,
		typeName:          typeName,
		isModel:           isModel,
		isMandatory:       isMandatory,
		isList:            isList,
		isMandatoryInList: isMandatoryInList,
		params:            params,
	}
	return a
}

func getModelLine(a attribute) string {
	if strings.EqualFold(a.name, "id") {
		a.name = strings.ToUpper(a.name)
	}

	if !a.isModel {
		switch a.typeName {
		case "String", "string":
			a.typeName = "string"
		case "Int", "int":
			a.typeName = "int32"
		case "Float", "float":
			a.typeName = "float64"
		case "Boolean", "boolean", "Bool", "bool":
			a.typeName = "bool"
		case "ID", "id":
			a.typeName = "graphql.ID"
		case "time", "Time":
			a.typeName = "time.Time"
		default:
			a.typeName = fmt.Sprintf("scalar.%s", strings.Title(a.typeName))
		}
	}

	if a.isList {
		if a.isMandatoryInList {
			a.typeName = "[]" + a.typeName
		} else {
			a.typeName = "[]*" + a.typeName
		}
	}

	if !a.isMandatory {
		a.typeName = "*" + a.typeName
	}

	if json {
		if a.isModel {
			return fmt.Sprintf("	%s %s `json:\"-\"`\n", strings.Title(a.name), strings.Title(a.typeName))
		} else if a.typeName == "*bool" || a.typeName == "bool" {
			return fmt.Sprintf("	%s %s `json:\"%s\"`\n", strings.Title(a.name), a.typeName, toSnake(a.name))
		} else {
			return fmt.Sprintf("	%s %s `json:\"%s,omitempty\"`\n", strings.Title(a.name), a.typeName, toSnake(a.name))
		}
	}
	if a.isModel {
		return fmt.Sprintf("	%s %s\n", strings.Title(a.name), strings.Title(a.typeName))
	}
	return fmt.Sprintf("	%s %s\n", strings.Title(a.name), a.typeName)
}

func getSchemaLine(a attribute) string {
	var typeReturn string
	switch a.typeName {
	case "boolean", "Bool", "bool":
		typeReturn = "Boolean"
	case "id", "Id":
		typeReturn = "ID"
	case "time", "Time":
		typeReturn = "String"
		if len(a.params) == 0 {
			p := attribute{name: "format", typeName: "String"}
			a.params = append(a.params, p)
		}
	default:
		typeReturn = strings.Title(a.typeName)
	}

	if a.isList {
		if a.isMandatoryInList {
			typeReturn = "[" + typeReturn + "!]"
		} else {
			typeReturn = "[" + typeReturn + "]"
		}
	}

	if a.isMandatory {
		typeReturn += "!"
	}

	if len(a.params) == 0 {
		return fmt.Sprintf("	%s: %s\n", a.name, typeReturn)
	}

	params := ""
	for i := range a.params {
		if i > 0 {
			params += ", "
		}
		str := strings.Replace(getSchemaLine(a.params[i]), "\n", "", -1)
		params += strings.TrimSpace(str)
	}

	return fmt.Sprintf("	%s(%s): %s\n", a.name, params, typeReturn)
}

func getResolverLine(a attribute) string {
	if !a.isModel {
		var typeReturn string
		var withFormat bool
		switch a.typeName {
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
			if len(a.params) == 0 {
				withFormat = true
			}
		default:
			typeReturn = fmt.Sprintf("scalar.%s", strings.Title(a.typeName))
		}

		if a.isList {
			if a.isMandatoryInList {
				typeReturn = "[]" + typeReturn
			} else {
				typeReturn = "[]*" + typeReturn
			}
		}

		if !a.isMandatory {
			typeReturn = "*" + typeReturn
		}

		params := ""
		for i := range a.params {
			if i == 0 {
				params += "args struct {\n"
			}
			params += getResolverParams(a.params[i]) + "\n"
			if i == len(a.params)-1 {
				params += "}"
			}
		}

		if !withFormat {
			return fmt.Sprintf(`func (r *{{.resolver}}) %s(%s) %s {
	return r.{{.abbreviation}}.%s
}
`, strings.Title(a.name), params, typeReturn, strings.Title(a.name))
		}

		return fmt.Sprintf(`func (r *{{.resolver}}) %s(args struct{ Format *string }) %s {
	return r.{{.abbreviation}}.Get%s(args.Format)
}
`, strings.Title(a.name), typeReturn, strings.Title(a.name))
	}

	a.typeName = fmt.Sprintf("%s%sResolver", strings.ToLower(string(a.typeName[0])), a.typeName[1:])

	if a.isList {
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
		if a.isMandatory {
			pointer = ""
			address = ""
			check = ""
		}

		if a.isMandatoryInList {
			insideAddress = "&"
		}

		bal = strings.Replace(bal, "{{.check}}", check, -1)
		bal = strings.Replace(bal, "{{.insideAddress}}", insideAddress, -1)
		bal = strings.Replace(bal, "{{.attribute}}", strings.Title(a.name), -1)
		bal = strings.Replace(bal, "{{.typeName}}", a.typeName, -1)
		bal = strings.Replace(bal, "{{.pointer}}", pointer, -1)
		bal = strings.Replace(bal, "{{.address}}", address, -1)

		return bal
	}

	address := ""
	if a.isMandatory {
		address = "&"
	}

	params := ""
	for i := range a.params {
		if i == 0 {
			params += "args struct {\n"
		}
		params += getResolverParams(a.params[i]) + "\n"
		if i == len(a.params)-1 {
			params += "}"
		}
	}

	return fmt.Sprintf(`func (r *{{.resolver}}) %s(%s) *%s {
	return &%s{%sr.{{.abbreviation}}.%s}
}
`, strings.Title(a.name), params, a.typeName, a.typeName, address, strings.Title(a.name))
}

func getResolverParams(a attribute) string {
	var typeReturn string
	if !a.isModel {
		switch a.typeName {
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
			if len(a.params) == 0 {
				p := attribute{name: "format", typeName: "String"}
				a.params = append(a.params, p)
			}
		default:
			typeReturn = fmt.Sprintf("scalar.%s", strings.Title(a.typeName))
		}
	} else {
		typeReturn = fmt.Sprintf("model.%s", strings.Title(a.typeName))
	}

	if a.isList {
		if a.isMandatoryInList {
			typeReturn = "[]" + typeReturn
		} else {
			typeReturn = "[]*" + typeReturn
		}
	}

	if !a.isMandatory {
		typeReturn = "*" + typeReturn
	}
	return fmt.Sprintf("	%s %s", strings.Title(a.name), typeReturn)
}

func getModelMethods(a attribute) string {
	if !a.isModel && len(a.params) == 0 {
		switch a.typeName {
		case "time", "Time":
			data := map[string]interface{}{"attribute": a.name, "abbreviation": toAbbreviation(a.name), "Attribute": strings.Title(a.name), "type": "%s"}
			if !a.isMandatory {
				data["notMandatory"] = "*"
			}
			if a.isList {
				data["list"] = "[]"
				if !a.isMandatoryInList {
					data["notInList"] = "*"
				}
			}
			return executeTemplate(templates["getDate"], data)
		}
	}
	return ""
}

func separateLine(line string) map[string]string {
	r := regexp.MustCompile(`\A(?P<attribute>[a-z][a-zA-Z0-9]*)(?P<params>\(([a-z][a-zA-Z0-9]*\: ?(([a-zA-Z0-9]+!?)|(\[[a-zA-Z0-9]+!?\]!?))(, ?)?)+\))?\: ?(?P<return>([a-zA-Z0-9]+!?)|(\[[a-zA-Z0-9]+!?\]!?))\z`)

	match := r.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			result[name] = match[i]
		}
	}
	return result
}
