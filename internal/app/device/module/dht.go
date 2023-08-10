package module

import "sync/atomic"

type DHT struct {
	*module
}

func newDHT(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &DHT{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeDTH,
		},
	}
}
