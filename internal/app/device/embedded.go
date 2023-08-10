package device

import "sync/atomic"

type embedded struct {
	*device
}

func newEmbedded(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &embedded{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeEmbedded,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
