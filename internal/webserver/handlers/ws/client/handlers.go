package client

import (
	"encoding/json"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	validate "github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/validator"
	"reflect"
)

// rpcHandler structure which describes how handler should look like,
// it more than enough for any cases
type rpcHandler struct {
	fn  reflect.Value // handler function which would be called for the API endpoint
	arg reflect.Type  // argument of this function, it represents as specific request structure, also can be <nil>
}

/*
parseHandler brings all functions to the same interface.

Handler function design will be looks like:
[] - means that this argument is optional

	func handlerName([r requestStruct]) ([responseStruct,] error) {
		// handler body...
	}

Note: if *responseStruct is <nil>, we should get not <nil> error
Note: if *responseStruct is not <nil>, we should get <nil> error
Note: responseStruct can be reference type
*/
func parseHandler(fn interface{}) rpcHandler {
	// check handler design as it described above
	// if any check fails we panic
	h := reflect.TypeOf(fn)
	if h.Kind() != reflect.Func {
		panic("function expected")
	}

	// check function arguments
	if h.NumIn() > 1 {
		panic("expected at most one handler argument")
	}

	// check function return values
	errIface := reflect.TypeOf((*error)(nil)).Elem()
	switch n := h.NumOut(); n {
	case 1, 2:
		if !h.Out(n - 1).Implements(errIface) {
			panic("at least one of return value must implement error")
		}
	default:
		panic("expected one or two return values")
	}

	// define request structure if we have it
	var req reflect.Type
	if n := h.NumIn(); n == 1 {
		req = h.In(n - 1)
	}

	return rpcHandler{
		fn:  reflect.ValueOf(fn),
		arg: req,
	}
}

/*
call makes function call and returns handler's response to the client

Firstly it parses and initializes function parameters with which function would be called.
Then makes function call with it and return handler's response.
*/
func (h *rpcHandler) call(params json.RawMessage) (json.RawMessage, error) {
	var in []reflect.Value
	if h.arg == nil {
		if params != nil {
			return nil, errors.InvalidParams.New("method hasn't any params")
		}
	} else {
		instance := reflect.New(h.arg).Interface()
		if err := json.Unmarshal(params, &instance); err != nil {
			return nil, errors.ParsingErr.Wrap(err, "params unmarshalling")
		}

		if err := validate.Struct(instance); err != nil {
			return nil, errors.ValidationErr.Use(err)
		}

		in = []reflect.Value{reflect.ValueOf(reflect.ValueOf(instance).Elem().Interface())} // dereferencing
	}

	// parse return structure ([*responseStruct,] error)
	switch ret, n := h.fn.Call(in), h.fn.Type().NumOut(); {
	case n == 1:
		if !ret[n-1].IsNil() {
			return nil, ret[n-1].Interface().(error)
		}

		return nil, nil
	case n == 2:
		if !ret[n-1].IsNil() {
			return nil, ret[n-1].Interface().(error)
		}

		res, err := json.Marshal(ret[0].Interface())
		return res, errors.ParsingErr.Wrap(err, "result unmarshalling")
	default:
		panic("expected 1 or 2 return values")
	}
}

/*
unwrap makes initializing with zero values for reference data structures.

Note: it works only for structures, not for pointers to structures. [t]
Note: it works only for structures zero values, not for pointers to zero value structures. [v]

param: t must be reflect.Struct type for data structure to unwrap
param: v must be reflect.New(typ Type).Elem() [typ - reflect.Struct]
*/
func unwrap(t reflect.Type, v reflect.Value) {
	// filtering pointers to the structure at first nesting level, in this case unwrap has no effect.
	if t.Kind() != reflect.Struct {
		return
	}

	for n, count := 0, v.NumField(); n < count; n++ {
		fieldType := t.Field(n)
		fieldValue := v.Field(n)
		switch fieldType.Type.Kind() {
		case reflect.Map:
			fieldValue.Set(reflect.MakeMap(fieldType.Type))
		case reflect.Ptr:
			// dereference ptrFieldValue, for be possible to unwrap pointer type field
			ptrFieldValue := reflect.New(fieldType.Type.Elem())
			unwrap(fieldType.Type.Elem(), ptrFieldValue.Elem())
			fieldValue.Set(ptrFieldValue)
		case reflect.Slice:
			fieldValue.Set(reflect.MakeSlice(fieldType.Type, 0, 0))
		case reflect.Struct:
			unwrap(fieldType.Type, fieldValue)
		}
	}
}
