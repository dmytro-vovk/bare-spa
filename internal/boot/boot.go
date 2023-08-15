package boot

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/dmytro-vovk/bare-spa/internal/app/guest"
	"github.com/dmytro-vovk/bare-spa/internal/app/user"
	"github.com/dmytro-vovk/bare-spa/internal/boot/config"
	"github.com/dmytro-vovk/bare-spa/internal/webserver"
	"github.com/dmytro-vovk/bare-spa/internal/webserver/handlers/home"
	"github.com/dmytro-vovk/bare-spa/internal/webserver/handlers/ws"
	"github.com/dmytro-vovk/bare-spa/internal/webserver/handlers/ws/client"
	"github.com/dmytro-vovk/bare-spa/internal/webserver/router"
)

type boot struct {
	container
	configPath string
}

func New(configPath string) (*boot, error) {
	b := &boot{
		configPath: configPath,
	}

	if err := b.loadConfig(); err != nil {
		return nil, err
	}

	go b.shutdown()

	return b, nil
}

func (b *boot) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	log.Printf("Got %v, shutting down...", s)
	b.container.shutdown()
	os.Exit(0)
}

func (b *boot) Config() *config.Config {
	return b.Get("Config").(*config.Config)
}

func (b *boot) loadConfig() error {
	s, err := config.Load(b.configPath)
	if err != nil {
		return err
	}

	log.Printf("Config loaded from %s", b.configPath)

	b.Set("Config", s, nil)

	return nil
}

func (b *boot) GuestApp() *guest.Guest {
	const id = "Guest Application"
	if s, ok := b.Get(id).(*guest.Guest); ok {
		return s
	}

	s := guest.New(b.SessionManager())

	b.Set(id, s, nil)

	return s
}

func (b *boot) UserApp() *user.User {
	const id = "User Application"
	if s, ok := b.Get(id).(*user.User); ok {
		return s
	}

	var address string
	if b.Config().WebServer.TLS.Enabled {
		address = "https" + "://" + b.Config().WebServer.Domain
	} else {
		address = "http" + "://" + b.Config().WebServer.Domain
	}

	s := user.New(address, b.SessionManager())

	b.WSUserClient().AddAPI(s)

	b.Set(id, s, nil)

	return s
}

func (b *boot) SessionManager() *scs.SessionManager {
	const id = "Sessions"
	if s, ok := b.Get(id).(*scs.SessionManager); ok {
		return s
	}

	s := scs.New()
	s.Store = memstore.New()
	s.Lifetime = time.Hour * 24 * 7
	s.Cookie.Persist = true
	s.Cookie.Secure = true
	s.Cookie.Name = "Host-session"

	b.Set(id, s, nil)

	return s
}

func (b *boot) Handler() *home.Handler {
	const id = "Web handler"
	if s, ok := b.Get(id).(*home.Handler); ok {
		return s
	}

	s := home.New(b.SessionManager())

	b.Set(id, s, nil)

	return s
}

func (b *boot) WebRouter() http.Handler {
	const id = "Web Router"
	if s, ok := b.Get(id).(http.Handler); ok {
		return s
	}

	h := b.Handler()

	r := b.SessionManager().LoadAndSave(router.New(
		router.Route("/ws", b.WebsocketHandler().Handler),
		router.Route("/js", h.Scripts),
		router.Route("/js.map", h.Scripts),
		router.Route("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
		router.RoutePrefix("/css/", home.Styles.ServeHTTP),
		router.CatchAll(h.Handle),
	))

	b.Set(id, r, nil)

	return r
}

func (b *boot) Webserver() *webserver.Webserver {
	const id = "Web Server"
	if s, ok := b.Get(id).(*webserver.Webserver); ok {
		return s
	}

	cfg := b.Config().WebServer

	var server *webserver.Webserver

	if b.Config().WebServer.TLS.Enabled {
		server = webserver.NewTLS(net.JoinHostPort(cfg.Domain, "443"), b.WebRouter(), cfg.TLS.CertDir, cfg.TLS.HostNames)
	} else {
		server = webserver.New(cfg.Domain, b.WebRouter())
	}

	b.GuestApp()
	b.UserApp()

	b.Set(id, server, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Println("Error stopping web server:", err)
		}
	})

	return server
}

func (b *boot) WebsocketHandler() *ws.Handler {
	const id = "WS Handler"
	if s, ok := b.Get(id).(*ws.Handler); ok {
		return s
	}

	h := ws.NewHandler(map[string]*client.Client{
		"guest": b.WSGuestClient(),
		"user":  b.WSUserClient(),
	}, b.SessionManager())

	b.Set(id, h, nil)

	return h
}

func (b *boot) WSGuestClient() *client.Client {
	const id = "WS Guest Client"
	if s, ok := b.Get(id).(*client.Client); ok {
		return s
	}

	s := client.New(b.SessionManager()).AddAPI(b.GuestApp())

	b.Set(id, s, nil)

	return s
}

func (b *boot) WSUserClient() *client.Client {
	const id = "WS User Client"
	if s, ok := b.Get(id).(*client.Client); ok {
		return s
	}

	s := client.New(b.SessionManager())

	b.Set(id, s, nil)

	return s
}
