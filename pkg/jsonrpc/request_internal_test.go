package jsonrpc

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValid(t *testing.T) {
	testCases := []struct {
		name     string
		request  json.RawMessage
		hasError bool
	}{
		{
			name:     "request with all fields",
			request:  json.RawMessage(`{ "jsonrpc": "2.0", "id": 1, "method": "makeAction", "params": 123 }`),
			hasError: false,
		},
		{
			name:     "notification request",
			request:  json.RawMessage(`{ "jsonrpc": "2.0", "method": "subscribe", "params": "method" }`),
			hasError: false,
		},
		{
			name:     "notification request without params",
			request:  json.RawMessage(`{ "jsonrpc": "2.0", "method": "disconnect" }`),
			hasError: false,
		},
		{
			name:     "request with unsupported protocol version",
			request:  json.RawMessage(`{ "jsonrpc": "1.0", "id": null, "method": "subscribe", "params": "method" }`),
			hasError: true,
		},
		{
			name:     "request without method",
			request:  json.RawMessage(`{ "jsonrpc": "2.0", "id": 2, "params": "method" }`),
			hasError: true,
		},
		{
			name:     "request with method which reserved for rpc-internal methods and extensions",
			request:  json.RawMessage(`{ "jsonrpc": "2.0", "id": 3, "method": "rpc.method", "params": "method" }`),
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req Request
			if err := json.Unmarshal(tc.request, &req); err != nil {
				t.Fatal("unmarshalling error:", err)
			}

			err := req.Valid()
			if tc.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
