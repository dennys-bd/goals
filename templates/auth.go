package templates

func initAuth() {
	Templates["auth"] = `

// AuthResolver type for graphql
type AuthResolver struct{
	u model.User
}

// FillAuthStruct puts the user inside the AuthResolver
func (r *Resolver) FillAuthStruct(ctx context.Context) {
	// TODO: Use ctx to get auth variables and set the AuthStruct
	// E.G.
	// id := ctx.Value(lib.ContextKeyAuth)
	// DB.find(id, &r.u)
}
`
}
