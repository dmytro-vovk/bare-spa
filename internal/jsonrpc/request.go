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

func (r Request) Response(data json.RawMessage) Response {
	return Response{
		ID:      r.ID,
		Version: "2.0",
		Result:  data,
	}
}

func (r Request) ErrorResponse(err error) Response {
	return Response{
		ID:      r.ID,
		Version: "2.0",
		Error:   err.Error(),
	}
}
