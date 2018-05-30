package cmd

func initGql() {
	Templates["gqltypes"] = `package gqltype

// {{.Gqltypes}} for graphql definition
const {{.Gqltypes}} = "type {{.Gqltype}} {\n" +
	my{{.Gqltypes}} + "\n" +
	"}"

// TODO: Concatenate here your {{.gqltypes}} types comming from schema
// for User {{.Gqltypes}}: const myQueries = schema.User{{.Gqltypes}}
const my{{.Gqltypes}} = ""
`

	Templates["scalar"] = `package gqltype

// Scalars for graphql definition
// TODO: Write here your scalars types
const Scalars = ` + "`" + `
scalar Time
` + "`" + `
`

	Templates["schema"] = `package gqltype

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

	Templates["types"] = `package gqltype

// Types for graphql definition
// TODO: Concatenate here your types comming from schema
// e.g. User Types: const Types = schema.UserTypes
const Types = ""
`
}
