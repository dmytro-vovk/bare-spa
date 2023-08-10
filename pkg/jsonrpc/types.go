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
		ID      int             `json:"id"`               // Must be the same as the value of the id member in the request
		Version string          `json:"jsonrpc"`          // Must be exactly "2.0"
		Result  json.RawMessage `json:"result,omitempty"` // The value is determined by the method invoked on the server
		Error   *Error          `json:"error,omitempty"`  // Returned object when a rpc call encounters an error
	}

	Error struct {
		Code    int             `json:"code"`           // The error type that occurred. [-32768; -32000] are reserved
		Message string          `json:"message"`        // Should be limited to a concise single sentence
		Data    json.RawMessage `json:"data,omitempty"` // Contains additional information about the error
	}
)

func (r Request) String() string { return stringify(r) }

func (r Response) String() string { return stringify(r) }

func (e Error) String() string { return stringify(e) }

func stringify(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
