package device

import "sync/atomic"

type ups12 struct {
	*device
}

func newUPS12(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &ups12{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeUPS12,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
