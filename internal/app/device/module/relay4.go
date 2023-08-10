package module

import "sync/atomic"

type Relay4 struct {
	*module
}

func newRelay4(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &Infrared{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeRelay4,
		},
	}
}
