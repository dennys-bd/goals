package templates

func initResolver() {
	Templates["fullresolver"] = `package resolver

{{.importpath}}type {{.resolver}} struct {
	{{.abbreviation}} *model.{{.Name}}
}

{{.methods}}
`
}
