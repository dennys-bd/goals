package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/dennys-bd/goals/auth"
	"github.com/dennys-bd/goals/graphql"
	"github.com/dennys-bd/goals/graphql/relay"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
)

var registers []register

type Schema struct {
	name   string
	schema string
}

type register struct {
	endpoint string
	schema   Schema
	resolver interface{}
	opt      []graphql.SchemaOpt
}

type CoreOpts struct {
	port     string
	verbose  bool
	graphiql bool
}

// RegisterSchema register your schema to a resolver in an endpoint
func RegisterSchema(endpoint string, schema Schema, resolver interface{}, opt ...graphql.SchemaOpt) {
	if endpoint == "/" {
		endpoint = ""
	}

	r := register{
		endpoint: endpoint,
		schema:   schema,
		resolver: resolver,
	}

	registers = append(registers, r)
}

// RegisterPrivateSchema register your private schema to a resolver in an endpoint
//
// RegisterPrivateSchema only calls RegisterSchema, but you may to use it if you want
// to garantee that your resolver is a PrivateResolver and you have a closed Schema.
func RegisterPrivateSchema(endpoint string, schema Schema, resolver graphql.PrivateResolver, opt ...graphql.SchemaOpt) {
	RegisterSchema(endpoint, schema, resolver, opt...)
}

// Server is user to run your goals application,
// User It after registering yours schemas.
func Server(opts CoreOpts) {
	for _, reg := range registers {
		if opts.graphiql || os.Getenv("GOALS_ENV") != "production" {
			innerPage := fmt.Sprintf(page, reg.endpoint)
			http.Handle(reg.endpoint+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(innerPage))
			}))
		}

		schema := graphql.MustParseSchema(reg.schema.schema, reg.resolver)
		if res, ok := reg.resolver.(graphql.PrivateResolver); ok {
			http.Handle("/graphql"+reg.endpoint, auth.InjectAuthToContext(&relay.Handler{Schema: schema}, res.GetAuthHeaders()...))
		} else {
			http.Handle("/graphql"+reg.endpoint, &relay.Handler{Schema: schema})
		}
	}
	go printServers(opts)
	log.Fatal(http.ListenAndServe(":"+opts.port, nil))
}

func GetOpts() CoreOpts {
	arg := os.Args[1]
	arg2, _ := strconv.ParseBool(os.Args[2])
	return CoreOpts{port: arg, verbose: arg2}
}

func printServers(opts CoreOpts) {
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-\n")
	for _, reg := range registers {
		fmt.Printf("%s is registered at: http://localhost:%s/graphql%s\nWith the resolver: %s\nYou can visit it's GraphiQL page on http://localhost:%s%s\n",
			reg.schema.name, opts.port, reg.endpoint, reflect.TypeOf(reg.resolver).Elem(), opts.port, reg.endpoint)
		if opts.verbose {
			fmt.Printf("Check the schema:\n`%s`\n", reg.schema.schema)
		}
		fmt.Printf("-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-=x=-\n")
		time.Sleep(500 * time.Millisecond)
	}
	println("Your server is running, press ctrl+c to stop it.")
}

// MountSchema from params
func MountSchema(name, types, queries, mutations, subscriptions, scalars string) Schema {
	schemaDefinition := "schema {\n"
	if queries != "" {
		schemaDefinition += "	query: Query\n"
	} else {
		errs.Ex(fmt.Sprintf("Your query can't be empty. Error in your schema: \"%s\"\n", name))
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

	s := Schema{
		name:   name,
		schema: schema,
	}

	return s
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
