package cmd

import (
	"path/filepath"
	"strings"

	"github.com/dennys-bd/goals/core"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	"github.com/spf13/cobra"
)

var noModel bool
var modelName string
var resolverName string

var authCmd = &cobra.Command{
	Use:     "authorization [name]",
	Aliases: []string{"auth"},
	Short:   "Generate structure for authorization",
	Long: `Authorization (goals scaffold authorization or 
symply goals s a) will create a new resolver 
and schema (and their structures), to be registered 
as a schema that must be private.
Only with your authorization you can allow access.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			errs.Ex("You can pass only one argument")
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

func createAuthFiles(project core.Project) {

	basicTemplates()
	gqlTemplates()

	modelName = strings.Title(modelName)
	resolverName = strings.Title(resolverName)

	abbreviation := toAbbreviation(modelName)

	resData := map[string]interface{}{"importpath": project.ImportPath, "name": resolverName, "model": "User", "abbreviation": abbreviation}
	resScript := executeTemplate(templates["resolver"], resData)
	writeStringToFile(filepath.Join(project.ResolverPath(), strings.ToLower(resolverName)+"_resolver.go"), resScript)

	schData := map[string]interface{}{"importpath": project.ImportPath, "Name": resolverName, "name": strings.ToLower(resolverName)}
	schScript := executeTemplate(templates["schema"], schData)
	writeStringToFile(filepath.Join(project.SchemaPath(), strings.ToLower(resolverName)+"_schema.go"), schScript)

	if !noModel {
		writeModel(modelName, "	ID		 graphql.ID\n	Name	 string\n	Email	 string\n	Password string\n", project)
	}
}
