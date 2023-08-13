package client

import (
	"context"
	"encoding/json"
	"reflect"
)

type rpcHandler struct {
	receiver   *reflect.Value
	fn         reflect.Value // handler function which would be called for the API endpoint
	arg        reflect.Type  // argument of this function, it represents as specific request structure, also can be <nil>
	hasContext bool          // if the first argument is context.Context
}

var (
	errIface     = reflect.TypeOf((*error)(nil)).Elem()
	ctx          = context.Background()
	contextIface = reflect.TypeOf(&ctx).Elem()
)

func parseMethod(fn reflect.Method, receiver *reflect.Value) *rpcHandler {
	var (
		arg        reflect.Type
		hasContext bool
	)

	if fn.Type.NumIn() > 3 {
		panic("expected at most two handler argument")
	} else if fn.Type.NumIn() == 3 {
		if !fn.Type.In(1).Implements(contextIface) {
			panic("context expected")
		}
		hasContext = true
		arg = fn.Type.In(2)
	} else if fn.Type.NumIn() == 2 {
		if fn.Type.In(1).Implements(contextIface) {
			hasContext = true
		} else {
			arg = fn.Type.In(1)
		}
	}

	switch n := fn.Type.NumOut(); n {
	case 1, 2:
		if !fn.Type.Out(n - 1).Implements(errIface) {
			panic("the last return value must implement error")
		}
	default:
		panic("expected one or two return values")
	}

	return &rpcHandler{
		receiver:   receiver,
		fn:         fn.Func,
		arg:        arg,
		hasContext: hasContext,
	}
}

func (h *rpcHandler) call(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	var ret, args []reflect.Value

	if h.receiver != nil {
		args = append(args, *h.receiver)
	}

	if h.hasContext {
		args = append(args, reflect.ValueOf(reflect.ValueOf(&ctx).Elem().Interface()))
	}

	if h.arg != nil {
		arg := reflect.New(h.arg).Interface()
		if err := json.Unmarshal(params, &arg); err != nil {
			return nil, err
		}

		args = append(args, reflect.ValueOf(reflect.ValueOf(arg).Elem().Interface()))
	}

	ret = h.fn.Call(args)

	switch n := h.fn.Type().NumOut(); {
	case n == 1:
		if ret[n-1].IsNil() {
			return nil, nil
		}

		return nil, ret[n-1].Interface().(error)
	case n == 2:
		if ret[n-1].IsNil() {
			return json.Marshal(ret[0].Interface())
		}

		return nil, ret[n-1].Interface().(error)
	default:
		panic("expected 1 or 2 return values")
	}
}
