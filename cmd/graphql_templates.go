package cmd

func gqlTemplates() {
	templates["scafmodel"] = `package model
	
{{.importpath}}// {{.Name}} Model
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
}
