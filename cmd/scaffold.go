package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:     "scaffold [name] atribute:type!",
	Aliases: []string{"scaf", "s"},
	Short:   "Creates attribute new object",
	Long: `Scaffold (goals scaffold) will create attribute new object
	and it's structure, nammed: Model. Schema and Resolvers.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Getwd()
		if err != nil {
			er(err)
		}

		_, err = ioutil.ReadFile("lib/goalsfile")
		if err != nil {
			er("This is not attribute goals project")
		}
	},
}

func writeAtribute(name string, args []string) (model string, schema string, resolver string) {

	var mB, sB, rB bytes.Buffer

	for _, attribute := range args {
		arguments := strings.Split(attribute, ":")
		if len(arguments) > 1 && len(arguments) < 4 {
			if arguments[1] == "type" && len(arguments) == 3 {
				mB.WriteString(getModelLine(arguments[0], arguments[2], true))
				sB.WriteString(getSchemaLine(arguments[0], arguments[2]))
				rB.WriteString(getResolverLine(arguments[0], arguments[2], true))

			} else {
				mB.WriteString(getModelLine(arguments[0], arguments[1], false))
				sB.WriteString(getSchemaLine(arguments[0], arguments[1]))
				rB.WriteString(getResolverLine(arguments[0], arguments[1], false))
			}
		} else {
			er("Bad Syntax")
		}
	}

	return model, schema, resolver
}

func getModelLine(attribute string, typeName string, isModel bool) string {
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
		return fmt.Sprintf("%s %s `json:\"-\"`\n", strings.Title(attribute), strings.Title(typeName))
	} else if typeName == "*bool" || typeName == "bool" {
		return fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	} else {
		return fmt.Sprintf("%s %s `json:\"%s,omitempty\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	}
}

func getSchemaLine(attribute string, typeName string) string {
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

	return fmt.Sprintf("%s: %s\n", attribute, typeName)
}

func getResolverLine(attribute string, typeName string, isModel bool) string {
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
			typeName = fmt.Sprintf("scalar.%s", typeName)
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

	typeName = fmt.Sprintf("*%sResolver", strings.Title(typeName))

	if isList {
		pointer := "*"
		address := "&"
		insideAddress := ""
		check := `if r.{{.abbreviation}}.{{.attribute}} == nil {
return nil
}`
		bal := `func (r *{{.resolver}}) {{.attribute}}() {{.pointer}}[]{{.typeName}} {
	{{.check}}
	slice := {{.pointer}}r.{{.abbreviation}}.{{.attribute}}

	l := make([]{{.typeName}}, len(slice))
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

		strings.Replace(bal, "{{.check}}", check, -1)
		strings.Replace(bal, "{{.insideAddress}}", insideAddress, -1)
		strings.Replace(bal, "{{.attribute}}", attribute, -1)
		strings.Replace(bal, "{{.typeName}}", typeName, -1)
		strings.Replace(bal, "{{.pointer}}", pointer, -1)
		strings.Replace(bal, "{{.address}}", address, -1)

	}

	if isModel {
		return fmt.Sprintf("%s %s `json:\"-\"`\n", strings.Title(attribute), typeName)
	} else if typeName == "*bool" || typeName == "bool" {
		return fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	} else {
		return fmt.Sprintf("%s %s `json:\"%s,omitempty\"`\n", strings.Title(attribute), typeName, toSnake(attribute))
	}
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
