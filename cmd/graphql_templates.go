package cmd

func gqlTemplates() {
	templates["scafmodel"] = `package model	

{{if .importpath}}{{.importpath}}{{end}}// {{.Name}} Model
type {{.Name}} struct {
{{.model}}}	
`

	templates["scafresolver"] = `package resolver

{{.importpath}}type {{.resolver}} struct {
	{{.abbreviation}} *model.{{.Name}}
}

{{.methods}} 
`

	templates["scafschema"] = `package schema

// {{.Name}}Types defines the graphql Types for {{.Name}}
// TODO: Concatenate {{.Name}}Types in gqltype.Types 
const {{.Name}}Types = ` + "`" + `
# {{.Name}} definition type
type {{.Name}} {
{{.schema}}}
` + "`" + `

// {{.Name}}Queries defines the graphql Queries for {{.Name}}
const {{.Name}}Queries = ""

// {{.Name}}Mutations defines the graphql Mutations for {{.Name}}
const {{.Name}}Mutations = ""
`

	templates["getDate"] = `
// Get{{.Attribute}} return model_name's {{.Attribute}} formated or not
func ({{.abbreviation}} model_name) Get{{.Attribute}}(format *string) {{.notMandatory}}{{.list}}{{.notInList}}string {
	{{if .list}}{{if .notMandatory}}if {{.abbreviation}}.{{.Attribute}} == nil {
		return nil
	}
	{{end}}dates := make([]{{.notInList}}string, len({{.notMandatory}}{{.abbreviation}}.{{.Attribute}}))
	for i := range dates {
		{{if .notInList}}dates[i] = getDateInFormat({{if .notMandatory}}(*{{.abbreviation}}.{{.Attribute}}){{else}}{{.abbreviation}}.{{.Attribute}}{{end}}[i], format)
	}{{else}}dates[i] = *(getDateInFormat(&{{if .notMandatory}}(*{{.abbreviation}}.{{.Attribute}}){{else}}{{.abbreviation}}.{{.Attribute}}{{end}}[i], format))
	}{{end}}
	return {{if .notMandatory}}&{{end}}dates
}{{else}}{{if .notMandatory}}return getDateInFormat({{.abbreviation}}.{{.Attribute}}, format){{else}}return *(getDateInFormat(&{{.abbreviation}}.{{.Attribute}}, format)){{end}}
}{{end}}
`
}
