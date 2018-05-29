package cmd

// Templates has all the templates to create Files
var Templates = make(map[string]string)

func init() {

	initServer()
	initConfig()
	initHelpers()
	initConfig()
	initGit()
}
