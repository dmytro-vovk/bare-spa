package module

import "sync/atomic"

type DS18B20 struct {
	*module
}

func newDS18B20(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &DS18B20{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeDS18B20,
		},
	}
}
