package client

import (
	ut "github.com/go-playground/universal-translator"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	methods     map[string]rpcHandler
	connections map[string]*connection
	mutex       sync.RWMutex
}

func New() *Client {
	return &Client{
		methods:     map[string]rpcHandler{},
		connections: map[string]*connection{},
	}
}

// NS adds handlers to the namespace
func (c *Client) NS(namespace string, handlers ...func(string, *Client)) *Client {
	for _, f := range handlers {
		f(namespace, c)
	}

	return c
}

// NSMethod add the handler to the namespace by name
func NSMethod(name string, handler interface{}) func(string, *Client) {
	return func(ns string, c *Client) {
		c.methods[ns+"."+name] = parseHandler(handler)
	}
}

// AddMethod add handler by name
func (c *Client) AddMethod(name string, fn interface{}) *Client {
	c.methods[name] = parseHandler(fn)
	return c
}

func (c *Client) Notify(method string, params interface{}) {
	if _, ok := params.(error); ok {
		panic("can't broadcast an error")
	}

	c.mutex.Lock()
	for addr := range c.connections {
		go c.connections[addr].Notify(method, params)
	}
	c.mutex.Unlock()
}

// Run handles single connection
func (c *Client) Run(conn *websocket.Conn, trans ut.Translator) {
	start, addr := time.Now(), conn.RemoteAddr().String()
	log.Printf("[%s] Websocket client connected", addr)
	c.connections[addr] = NewConnection(conn, c.methods, trans)
	c.connections[addr].Run()
	delete(c.connections, addr)
	log.Printf("[%s] Websocket client disconnected after %s", addr, time.Since(start))
}
