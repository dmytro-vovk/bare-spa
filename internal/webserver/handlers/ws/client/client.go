package client

import (
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
	"unicode"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/websocket"
)

type Client struct {
	session *scs.SessionManager
	methods map[string]*rpcHandler
	users   map[string]map[string]*connection // userID/address
	mutex   sync.RWMutex
}

func New(session *scs.SessionManager) *Client {
	return &Client{
		session: session,
		methods: map[string]*rpcHandler{},
		users:   make(map[string]map[string]*connection),
	}
}

func (c *Client) AddAPI(object any) *Client {
	ot, ov, namespace := reflect.TypeOf(object), reflect.ValueOf(object), ""
	if ot.Kind() == reflect.Pointer {
		namespace = canonical(ot.Elem().Name())
	} else {
		namespace = canonical(ot.Name())
	}

	for i := 0; i < ot.NumMethod(); i++ {
		m := ot.Method(i)
		name := namespace + "." + canonical(m.Name)
		c.methods[name] = parseMethod(m, &ov)
		log.Printf("Added RPC handler: %s", name)
	}

	return c
}

func (c *Client) NotifyUser(userID, method string, params any) {
	c.mutex.Lock()
	if conns, ok := c.users[userID]; ok {
		for i := range conns {
			go conns[i].Notify(method, params)
		}
	}
	c.mutex.Unlock()
}

func (c *Client) Run(r *http.Request, conn *websocket.Conn) {
	start, addr := time.Now(), conn.RemoteAddr().String()

	log.Printf("[%s] Websocket client connected", addr)

	cn := NewConnection(r, conn, c.methods)

	userID, ok := c.session.Get(r.Context(), "user_id").(string)

	c.mutex.Lock()
	if ok {
		if uc, ok := c.users[userID]; ok {
			uc[addr] = cn
		} else {
			c.users[userID] = map[string]*connection{addr: cn}
		}
	}
	c.mutex.Unlock()

	cn.Run()

	c.mutex.Lock()
	if ok {
		if uc, ok := c.users[userID]; ok {
			delete(uc, addr)

			if len(uc) == 0 {
				delete(c.users, userID)
			}
		}
	}
	c.mutex.Unlock()

	log.Printf("[%s] Websocket client disconnected after %s", addr, time.Since(start))
}

func canonical(name string) string {
	n := []rune(name)
	n[0] = unicode.ToLower(n[0])
	return string(n)
}
