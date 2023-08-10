package middleware

import (
	"github.com/Sergii-Kirichok/pr/internal/app/translator"
	ut "github.com/go-playground/universal-translator"
	"net/http"
	"strings"
)

// Translator just checking for query param & Accept-Language header but can be expanded to Cookie's etc....
func Translator() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var trans ut.Translator
			if locale := r.URL.Query().Get("locale"); len(locale) > 0 { // todo: it don't read it from URL, even if it exist
				var found bool
				if trans, found = translator.Get(locale); found {
					goto end
				}
			}

			// get and parse the "Accept-Language" HTTP header and return an array
			trans, _ = translator.Find(acceptedLanguages(r)...)
		end:
			r = r.WithContext(translator.WithContext(r.Context(), trans))
			next(w, r)
		}
	}
}

// acceptedLanguages returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func acceptedLanguages(r *http.Request) (languages []string) {
	accepted := r.Header.Get("Accept-Language")
	if accepted == "" {
		return
	}

	options := strings.Split(accepted, ",")
	l := len(options)
	languages = make([]string, l)
	for i := 0; i < l; i++ {
		locale := strings.SplitN(options[i], ";", 2)
		languages[i] = strings.Trim(locale[0], " ")
	}

	return
}
