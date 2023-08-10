package module

import "sync/atomic"

type Thermocouple struct {
	*module
}

func newThermocouple(id *uint64, config Config) Module {
	defer atomic.AddUint64(id, 1)
	return &Thermocouple{
		module: &module{
			id:   atomic.LoadUint64(id),
			name: config.Name,
			kind: TypeThermocouple,
		},
	}
}
