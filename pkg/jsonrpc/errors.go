package jsonrpc

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

const (
	// ParseError meaning: invalid JSON was received by the server. An error occurred on the server while parsing the JSON text
	ParseError = -32700
	// InvalidRequest meaning: the JSON sent is not a valid Request object
	InvalidRequest = -32600
	// MethodNotFound meaning: the method does not exist / is not available
	MethodNotFound = -32601
	// InvalidParams  meaning: invalid method parameter(s)
	InvalidParams = -32602
	// InternalError meaning: internal JSON-RPC error
	//InternalError = -32603

	// ServerErrorCodes - [-32000; -32099] meaning: reserved for implementation-defined server-errors.
	// The remainder of the space is available for application defined errors.
)

// Error returns the message of a jsonrpc.Error
func (e Error) Error() string {
	return e.Message
}

func NewError(code int, err error, data map[string]interface{}) *Error {
	if err == nil {
		return nil
	}

	var result []byte
	if data != nil {
		var e error
		result, e = json.Marshal(data)
		if e != nil {
			log.Panic(e)
		}
	}

	return &Error{
		Code:    code,
		Message: err.Error(),
		Data:    result,
	}
}
