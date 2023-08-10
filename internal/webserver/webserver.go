package webserver

import (
	"context"
	"errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/handlers/home"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/handlers/ws"
)

type Auth interface {
	Valid(username string, password string) bool
}

type Webserver struct {
	listen string
	server *http.Server
}

func New(handler *ws.Handler, listen string) *Webserver {
	return &Webserver{
		listen: listen,
		server: &http.Server{
			Addr: listen,
			Handler: NewRouter(
				Route("/ws", handler.Handler),
				Route("/js/index.js", home.Scripts),
				Route("/js/index.js.map", home.ScriptsMap),
				Route("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
					// TODO favicon.ico засунуть в html
					w.WriteHeader(http.StatusNoContent)
				}),
				RoutePrefix("/css/", home.Styles.ServeHTTP),
				CatchAll(home.Handler),
			),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (w *Webserver) Use(middlewares ...middleware.Middleware) *Webserver {
	l := len(middlewares)
	if l == 0 {
		return w
	}

	next := middlewares[l-1](w.server.Handler.ServeHTTP)
	for i := l - 2; i >= 0; i-- {
		next = middlewares[i](next)
	}

	w.server.Handler = next
	return w
}

func (w *Webserver) Serve() (err error) {
	log.Printf("Webserver starting at %s", w.listen)

	err = w.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return
}

func (w *Webserver) Stop(ctx context.Context) error { return w.server.Shutdown(ctx) }
