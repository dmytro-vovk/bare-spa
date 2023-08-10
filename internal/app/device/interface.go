package device

import "fmt"

type Interface string

const (
	InterfaceI2C   Interface = "I2C"
	InterfaceRS485 Interface = "RS485"
	InterfaceRadio Interface = "Radio"
	InterfaceSPI   Interface = "SPI"
	InterfaceUART  Interface = "UART"
)

func (kind Type) Interfaces() map[Interface]struct{} {
	return map[Type]map[Interface]struct{}{
		TypeEthIO2: {
			InterfaceRS485: {},
		},
		TypeEthIO4: {
			InterfaceRS485: {},
		},
		TypeIncubator1: {
			InterfaceRS485: {},
		},
		TypeIncubator2: {
			InterfaceRS485: {},
			InterfaceRadio: {},
		},
		TypeOmega1: {
			InterfaceUART: {},
			InterfaceSPI:  {},
		},
		TypeUPS12: {
			InterfaceRS485: {},
		},
		TypeEmbedded: {
			InterfaceUART: {},
		},
	}[kind]
}

type ErrInterfaceNotSupported struct {
	DeviceType      Type
	DeviceInterface Interface
}

func (e *ErrInterfaceNotSupported) Error() string {
	return fmt.Sprintf("interface %q not supported for device type %q", e.DeviceInterface, e.DeviceType)
}

func (kind Type) IsInterfaceSupported(iface Interface) bool {
	_, ok := kind.Interfaces()[iface]
	return ok
}
