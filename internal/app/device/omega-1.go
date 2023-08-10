package device

import "sync/atomic"

type omega1 struct {
	*device
}

func newOmega1(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &omega1{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeOmega1,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
