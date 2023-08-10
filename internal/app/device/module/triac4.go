package module

import "sync/atomic"

type Triac4 struct {
	*module
}

func newTriac4(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &Infrared{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeTriac4,
		},
	}
}
