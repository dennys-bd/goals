package templates

func initServer() {
	Templates["server"] = `package main
	
import (
	"log"
	"net/http"
	"os"

	"{{.importpath}}/app/gqltype"
	"{{.importpath}}/app/resolver"
	"{{.importpath}}/lib"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {

	if os.Getenv("ENVIRONMENT") != "production" {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(page)
		}))
	}

	// fmt.Println(gqltype.Schema)
	schema := graphql.MustParseSchema(gqltype.Schema, &resolver.Resolver{})

	http.Handle("/graphql", injectViewerToContext(&relay.Handler{Schema: schema}))

	port := lib.GetPort()

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func injectViewerToContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if (*r).Method == "OPTIONS" {
			return
		}

		// con := context.WithValue(r.Context(), lib.ContextKeyAuth, r.Header.Get("access-token"))
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

var page = []byte(` + "`" + `
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/graphql", {
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
` + "`" + `)
`
}
