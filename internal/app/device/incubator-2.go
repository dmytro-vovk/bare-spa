package device

import "sync/atomic"

type incubator2 struct {
	*device
}

func newIncubator2(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &incubator2{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeIncubator2,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
