package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "create"},
	Short:   "Initialize a Goals Application",
	Long: `Initialize (goals init) will create a new application, with 
the appropriate structure for a Go-Graphql application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Getwd()
		if err != nil {
			er(err)
		}

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
	if !exists(project.AbsPath()) {
		err := os.MkdirAll(project.AbsPath(), os.ModePerm)
		if err != nil {
			er(er)
		}
	} else if !isEmpty(project.AbsPath()) {
		er("Goals will not create a new project in a non empty directory: " + project.AbsPath())
	}

	initializeDep(project)
	createGQLTypes(project)
	createStructure(project)
	createAbsFiles(project)
	// createDatabase(project)
	// createAuth(project)
	downloadDepedences(project)

}

func initializeDep(project *Project) {
	cmd := exec.Command("dep", "init", project.AbsPath())
	err := cmd.Run()
	if err != nil {
		er(err)
	}
}

func createGQLTypes(project *Project) {
	queryData := map[string]interface{}{"Gqltypes": "Queries", "Gqltype": "Query", "gqltypes": "queries"}
	mutationData := map[string]interface{}{"Gqltypes": "Mutations", "Gqltype": "Mutation", "gqltypes": "mutations"}
	subscriptionData := map[string]interface{}{"Gqltypes": "Subscriptions", "Gqltype": "Subscription", "gqltypes": "subscriptions"}

	queryScript := executeTemplate(Templates["gqltypes"], queryData)
	mutationScript := executeTemplate(Templates["gqltypes"], mutationData)
	subscriptionScript := executeTemplate(Templates["gqltypes"], subscriptionData)

	writeStringToFile(filepath.Join(project.GqlPath(), "queries.go"), queryScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "mutations.go"), mutationScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "subscriptions.go"), subscriptionScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "scalars.go"), Templates["scalar"])
	writeStringToFile(filepath.Join(project.GqlPath(), "schema.go"), Templates["schema"])
	writeStringToFile(filepath.Join(project.GqlPath(), "types.go"), Templates["types"])
}

func createStructure(project *Project) {
	resolverTemplate := `package resolver

// Resolver type for graphql
type Resolver struct{}
`

	writeStringToFile(filepath.Join(project.ModelPath(), "helper.go"), Templates["modelHelper"])
	writeStringToFile(filepath.Join(project.ResolverPath(), "resolver.go"), resolverTemplate)
	writeStringToFile(filepath.Join(project.ResolverPath(), "helper.go"), Templates["resolverHelper"])

	err := os.MkdirAll(project.ScalarPath(), os.ModePerm)
	if err != nil {
		er(err)
	}

	err = os.MkdirAll(project.SchemaPath(), os.ModePerm)
	if err != nil {
		er(err)
	}
}

func createAbsFiles(project *Project) {
	serverData := map[string]interface{}{"importpath": project.ImportPath()}
	serverScript := executeTemplate(Templates["server"], serverData)

	procTemplate := `web: {{.appname}}`
	procData := map[string]interface{}{"appname": project.Name()}
	procScript := executeTemplate(procTemplate, procData)

	writeStringToFile(filepath.Join(project.AbsPath(), "server.go"), serverScript)
	writeStringToFile(filepath.Join(project.LibPath(), "config.go"), Templates["config"])
	writeStringToFile(filepath.Join(project.AbsPath(), ".gitignore"), Templates["git"])
	writeStringToFile(filepath.Join(project.AbsPath(), "Procfile"), procScript)
}

func downloadDepedences(project *Project) {
	cmd := exec.Command("dep", "ensure")
	cmd.Dir = project.AbsPath()
	err := cmd.Run()
	if err != nil {
		println(err)
		er(err)
	}
	// --add upper.io/db.v3
	// github.com/joho/godotenv
	// github.com/rs/cors
}
