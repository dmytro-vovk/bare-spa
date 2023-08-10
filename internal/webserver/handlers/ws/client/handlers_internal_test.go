package client

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testHandlerRequest struct {
	TestMap    map[int]bool
	TestSlice  []int
	TestPtr    *struct{ AnotherTestSlice []bool }
	TestStruct struct{ AnotherTestMap map[string]float64 }
}

func TestCall(t *testing.T) {
	type ret struct {
		msg json.RawMessage
		err error
	}

	testReq := testHandlerRequest{
		TestMap:   map[int]bool{0: false, 1: true},
		TestSlice: []int{0, 1},
		TestPtr: &struct{ AnotherTestSlice []bool }{
			AnotherTestSlice: []bool{true, false},
		},
		TestStruct: struct{ AnotherTestMap map[string]float64 }{
			AnotherTestMap: map[string]float64{"0": 0, "1": 1},
		},
	}

	testErr := errors.New("not <nil> error")

	testCases := []struct {
		name     string
		request  interface{}
		handler  interface{}
		expected ret
	}{
		{
			name:    "without request and with <nil> error",
			handler: func() error { return nil },
			expected: ret{
				msg: nil,
				err: nil,
			},
		},
		{
			name:    "without request and with not <nil> error",
			handler: func() error { return testErr },
			expected: ret{
				msg: nil,
				err: testErr,
			},
		},
		{
			name:    "with request and with not <nil> data and <nil> error in return statement",
			request: testReq,
			handler: func(thr testHandlerRequest) (*struct{}, error) { return &struct{}{}, nil },
			expected: ret{
				msg: json.RawMessage("{}"),
				err: nil,
			},
		},
		{
			name:    "with request and with <nil> data and not <nil> error in return statement",
			request: testReq,
			handler: func(thr testHandlerRequest) (*struct{}, error) { return nil, testErr },
			expected: ret{
				msg: nil,
				err: testErr,
			},
		},
		{
			name:    "data will be lost if you try to return them together with an error",
			request: testReq,
			handler: func(thr testHandlerRequest) (struct{}, error) { return struct{}{}, testErr },
			expected: ret{
				msg: nil,
				err: testErr,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rpc := parseHandler(tc.handler)

			var req []byte
			if tc.request != nil {
				var err error
				req, err = json.Marshal(tc.request)
				if err != nil {
					t.Fatalf("struct should be normally parsed")
				}
			}

			msg, err := rpc.call(req)
			//log.Printf("req: %q = %[1]x, err: %v", msg, err)

			assert.Equal(t, tc.expected.msg, msg)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUnwrap(t *testing.T) {
	testCases := []struct {
		name    string
		handler interface{}
	}{
		{
			name:    "regular request structure",
			handler: func(thr testHandlerRequest) error { return nil },
		},
		{
			name:    "reference request structure",
			handler: func(thr *testHandlerRequest) error { return nil },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rpc := parseHandler(tc.handler)
			unwrap(rpc.arg, reflect.New(rpc.arg).Elem())

			// todo: add checks
		})
	}
}
