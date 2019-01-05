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

	if os.Getenv("GOALS_ENV") != "production" {
		http.Handle("/", http.FileServer(http.Dir("./static")))
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
`
}
