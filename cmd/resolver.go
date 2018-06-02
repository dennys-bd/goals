package cmd

func initResolver() {
	Templates["fullresolver"] = `package resolver

import "{{.importpath}}/app/model"

type {{.resolver}} struct {
	{{.abbreviation}} *model.{{.Name}}
}

{{.methods}}
`
}
