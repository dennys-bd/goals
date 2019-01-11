package cmd

import (
	"path/filepath"
	"strings"

	"github.com/dennys-bd/goals/templates"
	"github.com/spf13/cobra"
)

var noModel bool
var modelName string
var resolverName string

var authCmd = &cobra.Command{
	Use:     "authorization [name]",
	Aliases: []string{"auth"},
	Short:   "Generate files for authorization",
	Long: `Authorization (goals auth) will create
a new resolver and it's structure to be parsed
with a private schema that must be private.
Only with your authorization you can allow access.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			er("You can pass only one argument")
		}

		project := recreateProjectFromGoals()

		createAuthFiles(project)
	},
}

func init() {
	authCmd.Flags().BoolVar(&noModel, "no-model", false, "Don't create a Model")
	authCmd.Flags().StringVarP(&modelName, "model", "m", "User", "Name to your auth model")
	authCmd.Flags().StringVarP(&resolverName, "resolver", "r", "Auth", "Name to your auth resolver")
}

func createAuthFiles(project Project) {
	modelName = strings.Title(modelName)
	resolverName = strings.Title(resolverName)

	abbreviation := toAbbreviation(modelName)

	resData := map[string]interface{}{"importpath": project.ImportPath, "name": resolverName, "model": "User", "abbreviation": abbreviation}
	resScript := executeTemplate(templates.Templates["resolver"], resData)
	writeStringToFile(filepath.Join(project.ResolverPath(), strings.ToLower(resolverName)+"_resolver.go"), resScript)

	schData := map[string]interface{}{"importpath": project.ImportPath, "Name": resolverName, "name": strings.ToLower(resolverName)}
	schScript := executeTemplate(templates.Templates["schema"], schData)
	writeStringToFile(filepath.Join(project.SchemaPath(), strings.ToLower(resolverName)+"_schema.go"), schScript)

	if !noModel {
		writeModel(modelName, "	ID		 graphql.ID\n	Name	 string\n	Email	 string\n	Password string\n", project)
	}
}
