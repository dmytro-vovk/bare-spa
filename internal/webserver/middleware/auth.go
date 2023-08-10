package middleware

import (
	"net/http"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func BasicAuth(auth interface{ Valid(string, string) bool }) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		if auth.Valid("", "") {
			return func(writer http.ResponseWriter, request *http.Request) {
				next(writer, request)
			}
		}

		return func(writer http.ResponseWriter, request *http.Request) {
			if username, password, ok := request.BasicAuth(); !ok {
				writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				writer.WriteHeader(http.StatusUnauthorized)
				return
			} else if !auth.Valid(username, password) {
				writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			next(writer, request)
		}
	}
}
