package device

import (
	"sync/atomic"
)

type ethIO2 struct {
	*device
}

func newEthIO2(id *uint64, config Config) Device {
	defer atomic.AddUint64(id, 1)
	return &ethIO2{
		device: &device{
			id:    atomic.LoadUint64(id),
			name:  config.Name,
			kind:  TypeEthIO2,
			iface: config.Interface,
			addr:  config.Address,
		},
	}
}
