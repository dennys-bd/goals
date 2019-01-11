package templates

func initServer() {
	Templates["server"] = `package main
	
import (
	"log"
	"net/http"

	"{{.importpath}}/app/resolver"
	"{{.importpath}}/app/schema"
	"{{.importpath}}/lib"

	"github.com/dennys-bd/goals/core"
)

func main() {

	core.StartWithResolver("/public", schema.GetSchema(), &resolver.Resolver{})

	port := lib.GetPort()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
`

	// func injectViewerToContext(next http.Handler) http.Handler {

	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// 		if (*r).Method == "OPTIONS" {
	// 			return
	// 		}

	// 		// con := context.WithValue(r.Context(), lib.ContextKeyAuth, r.Header.Get("access-token"))
	// 		next.ServeHTTP(w, r.WithContext(r.Context()))
	// 	})
	// }
}
