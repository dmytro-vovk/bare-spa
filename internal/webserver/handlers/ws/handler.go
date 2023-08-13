package ws

import (
	"log"
	"net/http"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/webserver/handlers/ws/client"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/websocket"
)

type Handler struct {
	roleClients map[string]*client.Client
	session     *scs.SessionManager
}

func NewHandler(c map[string]*client.Client, session *scs.SessionManager) *Handler {
	return &Handler{
		roleClients: c,
		session:     session,
	}
}

// Handler handles the websockets
func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	role, ok := h.session.Get(r.Context(), "role").(string)

	if !ok {
		role = "guest"
		h.session.Put(r.Context(), "role", role)
	}

	c, ok := h.roleClients[role]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := (&websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			_, ok := h.session.Get(r.Context(), "role").(string)
			return ok
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to websocket: %s", err)
		return
	}

	c.Run(r, conn)
}
