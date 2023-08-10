package module

import "sync/atomic"

type Resistive struct {
	*module
}

func newResistive(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &Resistive{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeResistive,
		},
	}
}
