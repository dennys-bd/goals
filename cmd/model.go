package cmd

func initModel() {
	Templates["fullmodel"] = `package model
	
// {{.Name}} Model
type {{.Name}} struct {
{{.model}}}	
`
}
