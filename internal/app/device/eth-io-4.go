package device

import "sync/atomic"

type ethIO4 struct {
	*device
}

func newEthIO4(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &ethIO4{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeEthIO4,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
