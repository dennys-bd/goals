package cmd

func basicTemplates() {
	templates["scalar"] = `package scalar

// Scalars for graphql definition
const Scalars = ` + "`" + `
scalar Time
` + "`" + `
`

	templates["schema"] = `package schema
	
import (
	"github.com/dennys-bd/goals/core"
	"{{.importpath}}/app/scalar"
)

// Concat here your types, queries, mutations and
// subscriptions that will be in the schema, e.g.
// const queries = userQueries +
// shoppingQueries

const {{if .name}}{{.name}}T{{else}}t{{end}}ypes = ""

const {{if .name}}{{.name}}Q{{else}}q{{end}}ueries = ""

const {{if .name}}{{.name}}M{{else}}m{{end}}utations = ""

const {{if .name}}{{.name}}S{{else}}s{{end}}ubscriptions = ""

// Get{{.Name}}Schema returns the schema String
func Get{{.Name}}Schema() string {
	{{if .name}}return core.MountSchema({{.name}}Types, {{.name}}Queries, {{.name}}Mutations, {{.name}}Subscriptions, scalar.Scalars)
	{{else}}return core.MountSchema(types, queries, mutations, subscriptions, scalar.Scalars){{end}}
}
`

	templates["resolver"] = `package resolver
{{if .importpath}}
import (
	"context"

	"{{.importpath}}/app/model"
	"{{.importpath}}/lib"
)
{{end}}
// {{if .name}}{{.name}}{{else}}Public{{end}}Resolver type for graphql
type {{if .name}}{{.name}}{{else}}Public{{end}}Resolver {{if .model}}struct { 
	{{.abbreviation}} model.{{.model}}
}{{else}}struct{}{{end}}
{{if .name}}
// FillAuthStruct puts the user inside the AuthResolver
func (r *{{.name}}Resolver) FillAuthStruct(ctx context.Context) {
	// TODO: Use ctx to get auth variables and set the AuthStruct
	// E.G.
	// id := ctx.Value(lib.ContextKeyAuth)
	// DB.find(id, &r.u)
}

// GetAuthHeaders put in the context the headers you want
// from the original request to the resolver
func (r *{{.name}}Resolver) GetAuthHeaders() []string {
	// TODO: return here the headers you want
	// to come in context from the original request
	return []string{lib.ContextKeyAuth.String()}
}
{{end}}`
}
