package device

import "sync/atomic"

type incubator1 struct {
	*device
}

func newIncubator1(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &incubator1{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeIncubator1,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
