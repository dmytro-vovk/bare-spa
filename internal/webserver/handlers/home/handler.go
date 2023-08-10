package home

import (
	"embed"
	_ "embed"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:embed index.html.gz
var homePage []byte

// Handler serves the all-in-one home page
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("From %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write(homePage); err != nil {
		log.Printf("Error writing response to %s: %s", r.RemoteAddr, err)
	}
}

//go:embed index.js.gz
var scripts []byte

func Scripts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/javascript")

	if _, err := w.Write(scripts); err != nil {
		log.Printf("Error writing response to %s: %s", r.RemoteAddr, err)
	}
}

//go:embed index.js.map.gz
var scriptsMap []byte

func ScriptsMap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/javascript")

	if _, err := w.Write(scriptsMap); err != nil {
		log.Printf("Error writing response to %s: %s", r.RemoteAddr, err)
	}
}

//go:embed css
var styles embed.FS

var Styles = http.FileServer(http.FS(styles))
