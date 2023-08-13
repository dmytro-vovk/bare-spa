package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/jsonrpc"
	"github.com/gorilla/websocket"
)

type connection struct {
	request       *http.Request
	conn          *websocket.Conn
	methods       map[string]*rpcHandler
	subscriptions map[string]struct{}
	sendC         chan any
	doneC         chan struct{}
	mutex         sync.RWMutex
}

func NewConnection(request *http.Request, conn *websocket.Conn, methods map[string]*rpcHandler) *connection {
	return &connection{
		request:       request,
		conn:          conn,
		methods:       methods,
		subscriptions: map[string]struct{}{},
		sendC:         make(chan any, 10),
		doneC:         make(chan struct{}),
	}
}

func (c *connection) Run() {
	go c.reader()
	go c.sender()

	<-c.doneC
}

func (c *connection) Notify(method string, params any) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if _, ok := c.subscriptions[method]; !ok {
		return
	}

	payload, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	c.notify(jsonrpc.Request{
		Version: "2.0",
		Method:  method,
		Params:  payload,
	})
}

func (c *connection) notify(notice jsonrpc.Request) bool {
	select {
	case c.sendC <- notice:
		return true
	default:
		log.Printf("[%s] Couldn't send notification", c.conn.RemoteAddr())
		return false
	}
}

func (c *connection) reader() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[%s] Unexpected close error: %v", c.conn.RemoteAddr(), err)
			}

			close(c.doneC)
			return
		}

		switch msgType {
		case websocket.TextMessage:
			go c.handleTextMessage(msg)
		default:
			log.Printf("[%s] Unknown message type: %d", c.conn.RemoteAddr(), msgType)
		}
	}
}

func (c *connection) sender() {
	for {
		select {
		case resp := <-c.sendC:
			switch t := resp.(type) {
			case jsonrpc.Response:
			case jsonrpc.Request:
			default:
				log.Panicf("unknown response type: %T", t)
			}

			if err := c.conn.WriteJSON(resp); err != nil {
				log.Printf("[%s] Error sending message: %s", c.conn.RemoteAddr(), err)
			}
		case <-c.doneC:
			return
		}
	}
}

func (c *connection) handleTextMessage(msg []byte) {
	var req jsonrpc.Request
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Printf("[%s] Error decoding request: %s", c.conn.RemoteAddr(), err)
		log.Printf("[%s] Request: %s", c.conn.RemoteAddr(), msg)
		c.sendC <- req.ErrorResponse(err)
		return
	}

	if err := req.Valid(); err != nil {
		log.Printf("[%s] Invalid request object: %s", c.conn.RemoteAddr(), err)
		c.sendC <- req.ErrorResponse(err)
		return
	}

	c.handleRequest(req)
}

func (c *connection) handleRequest(req jsonrpc.Request) {
	if req.IsNotification() {
		c.handleNotification(req)
		return
	}

	if fn, ok := c.methods[req.Method]; ok {
		if data, err := fn.call(context.WithValue(context.Background(), "request", c.request), req.Params); err != nil {
			log.Printf("[%s] RPC call %s(%s) error: %s", c.conn.RemoteAddr(), req.Method, req.Params, err.Error())
			c.sendC <- req.ErrorResponse(err)
		} else {
			c.sendC <- req.Response(data)
		}

		return
	}

	log.Printf("[%s] Requested method %q doesn't exist", c.conn.RemoteAddr(), req.Method)
	c.sendC <- req.ErrorResponse(fmt.Errorf("method %q doesn't exist", req.Method))
}

func (c *connection) handleNotification(notice jsonrpc.Request) {
	var method string
	if err := json.Unmarshal(notice.Params, &method); err != nil {
		log.Printf("[%s] Error decoding method name: %s", c.conn.RemoteAddr(), err)
		log.Printf("[%s] Params: %s", c.conn.RemoteAddr(), notice.Params)
		return
	}

	switch notice.Method {
	case "subscribe":
		c.subscribe(method)
	case "unsubscribe":
		c.unsubscribe(method)
	}
}

func (c *connection) subscribe(method string) {
	log.Printf("[%s] Subscribing to %q", c.conn.RemoteAddr(), method)
	c.mutex.Lock()
	c.subscriptions[method] = struct{}{}
	c.mutex.Unlock()
}

func (c *connection) unsubscribe(method string) {
	log.Printf("[%s] Unsubscribing from %q", c.conn.RemoteAddr(), method)
	c.mutex.Lock()
	delete(c.subscriptions, method)
	c.mutex.Unlock()
}
