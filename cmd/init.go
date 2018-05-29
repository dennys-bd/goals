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

	createGQLTypes(project)
	createModel(project)
	createResolver(project)
	createFolders(project)
	createServerFile(project)
	createConfigFile(project)
	createProcFile(project)
	createGitFile(project)
	// createDatabase(project)
	downloadDepedences(project)

}

func createGQLTypes(project *Project) {
	cmd := exec.Command("dep", "init", project.AbsPath())
	err := cmd.Run()
	if err != nil {

	}
	gqlTemplate := `package gqltype

// {{.Gqltypes}} for graphql definition
const {{.Gqltypes}} = "type {{.Gqltype}} {\n" +
	my{{.Gqltypes}} + "\n" +
	"}"

// TODO: Concatenate here your {{.gqltypes}} types comming from schema
// for User {{.Gqltypes}}: const myQueries = schema.User{{.Gqltypes}}
const my{{.Gqltypes}} = ""
`

	queryData := map[string]interface{}{"Gqltypes": "Queries", "Gqltype": "Query", "gqltypes": "queries"}
	mutationData := map[string]interface{}{"Gqltypes": "Mutations", "Gqltype": "Mutation", "gqltypes": "mutations"}
	subscriptionData := map[string]interface{}{"Gqltypes": "Subscriptions", "Gqltype": "Subscription", "gqltypes": "subscriptions"}

	queryScript := executeTemplate(gqlTemplate, queryData)
	mutationScript := executeTemplate(gqlTemplate, mutationData)
	subscriptionScript := executeTemplate(gqlTemplate, subscriptionData)

	writeStringToFile(filepath.Join(project.GqlPath(), "queries.go"), queryScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "mutations.go"), mutationScript)
	writeStringToFile(filepath.Join(project.GqlPath(), "subscriptions.go"), subscriptionScript)

	scalarTemplate := `package gqltype

// Scalars for graphql definition
// TODO: Write here your scalars types
const Scalars = ` + "`" + `
scalar Time
` + "`" + `
`
	writeStringToFile(filepath.Join(project.GqlPath(), "scalars.go"), scalarTemplate)

	schemaTemplate := `package gqltype

// TODO: Uncomment here what you will use
const schemaDefinition = ` + "`" + `
schema {
	#query: Query
	#mutation: Mutation
	#subscription: Subscription
}
` + "`" + `

// Schema concatenated
const Schema = schemaDefinition +
	Queries +
	Mutations +
	Scalars +
	Types
`

	writeStringToFile(filepath.Join(project.GqlPath(), "schema.go"), schemaTemplate)

	typesTemplate := `package gqltype

// Types for graphql definition
// TODO: Concatenate here your types comming from schema
// e.g. User Types: const Types = schema.UserTypes
const Types = ""
`

	writeStringToFile(filepath.Join(project.GqlPath(), "types.go"), typesTemplate)
}

func createModel(project *Project) {
	writeStringToFile(filepath.Join(project.ModelPath(), "helper.go"), Templates["modelHelper"])
}

func createResolver(project *Project) {
	resolverTemplate := `package resolver

// Resolver type for graphql
type Resolver struct{}
`

	writeStringToFile(filepath.Join(project.ResolverPath(), "resolver.go"), resolverTemplate)
	writeStringToFile(filepath.Join(project.ResolverPath(), "helper.go"), Templates["resolverHelper"])
}

func createFolders(project *Project) {
	err := os.MkdirAll(project.ScalarPath(), os.ModePerm)
	if err != nil {
		er(err)
	}

	err = os.MkdirAll(project.SchemaPath(), os.ModePerm)
	if err != nil {
		er(err)
	}
}

func createServerFile(project *Project) {
	data := map[string]interface{}{"importpath": project.ImportPath()}
	serverScript := executeTemplate(Templates["server"], data)

	writeStringToFile(filepath.Join(project.AbsPath(), "server.go"), serverScript)
}

func createConfigFile(project *Project) {
	writeStringToFile(filepath.Join(project.LibPath(), "config.go"), Templates["config"])
}

func createProcFile(project *Project) {
	template := `web: {{.appname}}`
	data := map[string]interface{}{"appname": project.Name()}
	script := executeTemplate(template, data)
	writeStringToFile(filepath.Join(project.AbsPath(), "Procfile"), script)
}

func createGitFile(project *Project) {
	writeStringToFile(filepath.Join(project.AbsPath(), ".gitignore"), Templates["git"])
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
