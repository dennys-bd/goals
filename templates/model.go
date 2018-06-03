package templates

func initModel() {
	Templates["fullmodel"] = `package model
	
{{.importpath}}// {{.Name}} Model
type {{.Name}} struct {
{{.model}}}	
`
}
