package core

import (
	"context"
	"net/http"
)

func injectAuthToContext(next http.Handler, headers ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}
		con := r.Context()
		for _, h := range headers {
			s := r.Header.Get(h)
			con = context.WithValue(con, h, s)
		}
		next.ServeHTTP(w, r.WithContext(con))
	})
}
