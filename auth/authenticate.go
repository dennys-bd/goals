package auth

import (
	"context"
	"net/http"

	"github.com/dennys-bd/letest/lib"
)

func injectViewerToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}

		con := context.WithValue(r.Context(), lib.ContextKeyAuth, viewerID)
		next.ServeHTTP(w, r.WithContext(con))
	})
}
