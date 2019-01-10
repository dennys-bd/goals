package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dennys-bd/goals/templates"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init name",
	Aliases: []string{"initialize", "create"},
	Short:   "Initialize a Goals Application",
	Long: `Initialize (goals init) will create a new application, with 
the appropriate structure for a Go-Graphql application.

* If a name or relative path is provided, it will be created inside $GOPATH;
  (e.g. github.com/dennys-bd/goals);
* If an absolute path is provided, it will be created INSIDE $GOPATH;
* If your working directory is inside $GOPATH, it will be created on the wd;
* If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Getwd()
		check(err)

		var project *Project
		if len(args) == 0 {
			er("please insert the project name")
		} else if len(args) == 1 {
			project = NewProject(args[0])
		} else {
			er("please provide only one argument")
		}

		intializeProject(project)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		print("Creating your project")
		go func() {
			for i := 0; i < 2; i++ {
				time.Sleep(500 * time.Millisecond)
				print(".")
			}
			time.Sleep(500 * time.Millisecond)
			println(".")
		}()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		println("Done!")
	},
}

func intializeProject(project *Project) {
	if !exists(project.AbsPath) {
		err := os.MkdirAll(project.AbsPath, os.ModePerm)
		check(err)
	} else if !isEmpty(project.AbsPath) {
		er("Goals will not create a new project in a non empty directory: " + project.AbsPath)
	}

	initializeDep(project)
	createGQLTypes(project)
	createStructure(project)
	createAbsFiles(project)
	// TODO: createDatabase(project)
	// TODO: createDotEnv(project)
	downloadDepedences(project)

}

func initializeDep(project *Project) {
	cmd := exec.Command("dep", "init", project.AbsPath)
	err := cmd.Run()
	check(err)

	str := `package main

	func main(){
	}
`
	writeStringToFile(filepath.Join(project.AbsPath, "main.go"), str)

	repositories := []string{"github.com/dennys-bd/goals"}

	for _, s := range repositories {
		cmd := exec.Command("dep", "ensure", "-add", s)
		cmd.Dir = project.AbsPath
		out, err := cmd.Output()
		check(err)
		println(string(out))
	}

	removeFile(filepath.Join(project.AbsPath, "main.go"))
}

func createGQLTypes(project *Project) {
	queryData := map[string]interface{}{"Gqltypes": "Queries", "Gqltype": "Query", "gqltypes": "queries"}
	mutationData := map[string]interface{}{"Gqltypes": "Mutations", "Gqltype": "Mutation", "gqltypes": "mutations"}
	subscriptionData := map[string]interface{}{"Gqltypes": "Subscriptions", "Gqltype": "Subscription", "gqltypes": "subscriptions"}

	queryScript := executeTemplate(templates.Templates["gqltypes"], queryData)
	mutationScript := executeTemplate(templates.Templates["gqltypes"], mutationData)
	subscriptionScript := executeTemplate(templates.Templates["gqltypes"], subscriptionData)

	writeStringToFile(filepath.Join(project.GqlPath(), "queries.go"), queryScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "mutations.go"), mutationScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "subscriptions.go"), subscriptionScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "scalars.go"), templates.Templates["scalar"])
	writeStringToFile(filepath.Join(project.GqlPath(), "schema.go"), templates.Templates["schema"])
	writeStringToFile(filepath.Join(project.GqlPath(), "types.go"), templates.Templates["types"])
}

func createStructure(project *Project) {
	resolverTemplate := `package resolver

// Resolver type for graphql
type Resolver struct{}
`

	writeStringToFile(filepath.Join(project.ModelPath(), "helper.go"), templates.Templates["modelHelper"])
	writeStringToFile(filepath.Join(project.ResolverPath(), "resolver.go"), resolverTemplate)
	writeStringToFile(filepath.Join(project.ResolverPath(), "helper.go"), templates.Templates["resolverHelper"])

	err := os.MkdirAll(project.ScalarPath(), os.ModePerm)
	check(err)

	err = os.MkdirAll(project.SchemaPath(), os.ModePerm)
	check(err)

}

func createAbsFiles(project *Project) {
	serverData := map[string]interface{}{"importpath": project.ImportPath}
	serverScript := executeTemplate(templates.Templates["server"], serverData)

	procTemplate := `web: {{.appname}}`
	procData := map[string]interface{}{"appname": project.Name}
	procScript := executeTemplate(procTemplate, procData)

	writeStringToFile(filepath.Join(project.AbsPath, "server.go"), serverScript)
	writeStringToFile(filepath.Join(project.AbsPath, ".gitignore"), templates.Templates["git"])
	writeStringToFile(filepath.Join(project.AbsPath, "Procfile"), procScript)
	writeStringToFile(filepath.Join(project.LibPath(), "config.go"), templates.Templates["config"])
	writeStringToFile(filepath.Join(project.LibPath(), "Goals.toml"), project.CreateGoalsToml())
	writeStringToFile(filepath.Join(project.StaticPath(), "index.html"), templates.Templates["index.html"])
	// TODO: Create consts.go
}

func downloadDepedences(project *Project) {
	cmd := exec.Command("dep", "ensure")
	cmd.Dir = project.AbsPath
	err := cmd.Run()
	check(err)
}
