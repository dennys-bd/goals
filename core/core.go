package core

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dennys-bd/goals/graphql"
	"github.com/dennys-bd/goals/graphql/relay"
)

// StartWithResolver Starts the resolver's endpoint
func StartWithResolver(endpoint, schemaString string, resolver interface{}, opt ...graphql.SchemaOpt) {

	if endpoint == "/" {
		endpoint = ""
	}

	if os.Getenv("GOALS_ENV") != "production" {
		http.Handle(endpoint+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf(page, endpoint)))
		}))
	}

	schema := graphql.MustParseSchema(schemaString, resolver, opt...)
	http.Handle("/graphql"+endpoint, &relay.Handler{Schema: schema})
}

// MountSchema from params
func MountSchema(types, queries, mutations, subscriptions, scalars string) string {
	schemaDefinition := "schema {\n"
	if queries != "" {
		schemaDefinition += "	query: Query\n"
	}
	if mutations != "" {
		schemaDefinition += "	mutation: Mutation\n"
	}
	if subscriptions != "" {
		schemaDefinition += "	subscription: Subscription\n"
	}
	schemaDefinition += "}\n"

	schema := schemaDefinition
	if queries != "" {
		q := "type Query {" + queries + "}\n"
		schema += q
	}
	if mutations != "" {
		m := "type Mutation {" + mutations + "}\n"
		schema += m
	}
	if subscriptions != "" {
		s := "type Subscription {" + subscriptions + "}\n"
		schema += s
	}
	if scalars != "" {
		schema += scalars
	}
	if types != "" {
		schema += types
	}
	return schema
}

const page = `<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%%; height: 100%%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/graphql%s", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`
