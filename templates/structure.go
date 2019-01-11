package templates

func initStructure() {
	Templates["scalar"] = `package scalar

// Scalars for graphql definition
const Scalars = ` + "`" + `
scalar Time
` + "`" + `
`

	Templates["schema"] = `package schema
	
import (
	"github.com/dennys-bd/goals/core"
	"{{.importpath}}/app/scalar"
)

// Concat here your types, queries, mutations and subscriptions
// that will be in the general schema e.g.
// const queries = userQueries +
// shoppingQueries

const types = ""

const queries = ""

const mutations = ""

const subscriptions = ""

// GetSchema returns the schema String
func GetSchema() string {
	return core.MountSchema(types, queries, mutations, subscriptions, scalar.Scalars)
}
`

	Templates["resolver"] = `package resolver
{{if .imports}}
{{.imports}}
{{end}}
// {{.name}}Resolver type for graphql
type {{.name}}Resolver struct{ {{.model}} }
{{if .name}}
// FillAuthStruct bla
func (r *{{.name}}Resolver) FillAuthStruct(ctx context.Context) {
	// fmt.Printf("On Resolver: %v\n", ctx.Value(lib.ContextKeyAuth))
}
{{end}}`
}
