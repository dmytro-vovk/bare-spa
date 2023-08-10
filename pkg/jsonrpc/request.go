package jsonrpc

import (
	"encoding/json"
	"errors"
	"strings"
)

func (r Request) Valid() error {
	if r.Version != "2.0" {
		return errors.New("unsupported protocol version")
	}

	if r.Method == "" || strings.HasPrefix(r.Method, "rpc.") {
		return errors.New("rpc-reserved or empty method")
	}

	return nil
}

func (r Request) IsNotification() bool {
	return r.ID == 0
}

// Either the result member or error member must be included, but both members mustn't be included.
func (r Request) Response(data json.RawMessage, err *Error) Response {
	if data == nil && err == nil {
		data = json.RawMessage(`{ "ok": true }`)
	}

	return Response{
		ID:      r.ID,
		Version: "2.0",
		Result:  data,
		Error:   err,
	}
}
