package webserver

import (
	"net/http"
	"strings"
)

type router struct {
	routes   map[string]http.HandlerFunc
	prefixes map[string]http.HandlerFunc
	catchAll http.HandlerFunc
}

func Route(name string, handler http.HandlerFunc) func(r *router) {
	return func(r *router) {
		r.routes[name] = handler
	}
}

func RoutePrefix(prefix string, handler http.HandlerFunc) func(r *router) {
	return func(r *router) {
		r.prefixes[prefix] = handler
	}
}

func CatchAll(handler http.HandlerFunc) func(r *router) {
	return func(r *router) {
		r.catchAll = handler
	}
}

func NewRouter(routes ...func(r *router)) http.HandlerFunc {
	r := router{
		routes:   map[string]http.HandlerFunc{},
		prefixes: map[string]http.HandlerFunc{},
	}

	for i := 0; i < len(routes); i++ {
		routes[i](&r)
	}

	return r.matcher()
}

func (r *router) matcher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if handler, ok := r.routes[request.URL.Path]; ok {
			handler(writer, request)

			return
		}

		for prefix, handler := range r.prefixes {
			if strings.HasPrefix(request.URL.Path, prefix) {
				handler(writer, request)

				return
			}
		}

		if r.catchAll != nil {
			r.catchAll.ServeHTTP(writer, request)

			return
		}

		writer.WriteHeader(http.StatusNotFound)
	}
}
