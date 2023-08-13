package home

import (
	"embed"
	_ "embed"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	session *scs.SessionManager
}

func New(sessionManager *scs.SessionManager) *Handler {
	return &Handler{session: sessionManager}
}

//go:embed index.html.gz
var indexPage []byte

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	role, ok := h.session.Get(r.Context(), "role").(string)
	if !ok {
		role = "guest"
		h.session.Put(r.Context(), "role", role)
	}

	log.Printf("[%s] %s %s as %s", r.RemoteAddr, r.Method, r.URL.Path, role)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("Referrer-Policy", "same-origin")
	w.Header().Add("Permissions-Policy", "geolocation=(), microphone=(), camera=(), display-capture=()")
	w.Header().Add("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com:443 https://adminlte.io:443")

	_, _ = w.Write(indexPage)
}

//go:embed css
var styles embed.FS

var Styles = http.FileServer(http.FS(styles))

//go:embed guest.js.gz
var guestScripts []byte

//go:embed guest.js.map.gz
var guestScriptsMap []byte

//go:embed user.js.gz
var userScripts []byte

//go:embed user.js.map.gz
var userScriptsMap []byte

var roleScripts = map[string][]byte{
	"guest/js":     guestScripts,
	"guest/js.map": guestScriptsMap,
	"user/js":      userScripts,
	"user/js.map":  userScriptsMap,
}

func (h *Handler) Scripts(w http.ResponseWriter, r *http.Request) {
	role, ok := h.session.Get(r.Context(), "role").(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, ok := roleScripts[role+r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/javascript")

	_, _ = w.Write(data)
}
