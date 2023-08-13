package jsonrpc

import "encoding/json"

type (
	Request struct {
		ID      int             `json:"id,omitempty"`     // 0 value is reserved for notification
		Version string          `json:"jsonrpc"`          // Must be exactly "2.0"
		Method  string          `json:"method"`           // Name of the method to be invoked
		Params  json.RawMessage `json:"params,omitempty"` // Values to be used during the invocation of the method
	}

	Response struct {
		ID      int             `json:"id"`              // Must be the same as the value of the id member in the request
		Version string          `json:"jsonrpc"`         // Must be exactly "2.0"
		Result  json.RawMessage `json:"result"`          // The value is determined by the method invoked on the server
		Error   string          `json:"error,omitempty"` // Returned object when a rpc call encounters an error
	}
)
