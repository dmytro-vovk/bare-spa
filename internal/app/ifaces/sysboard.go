package ifaces

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device/address"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device/module"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/i2c"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/spi"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/uart"
)

type Sysboard interface {
	Type() string
	Setup() error
	GetInterfaces() []iface.Interface
	GetTemperatures() (map[string]int16, error)
	GetUARTConfig(id byte) (*uart.Config, error)
	SetUARTConfig(id byte, config uart.Config) error
	GetSPIConfig(id byte) (*spi.Config, error)
	SetSPIConfig(id byte, config spi.Config) error
	GetI2CConfig(id byte) (*i2c.Config, error)
	SetI2CConfig(id byte, config i2c.Config) error
	SupportedDevicesTypes() []device.Type
	SupportedDeviceInterfaces(kind device.Type) []device.Interface
	SupportedDeviceModules(kind device.Type) []module.Type
	AvailableDevicesAddresses() []address.Address
	ConnectedDevicesList(offset, limit int) ([]*device.Entity, error) // todo: mb real []device.Device?
	ConnectedDevicesCount() int
	CreateDevice(kind device.Type, config device.Config) (device.Device, error)
	ReadDevice(id uint64) (device.Device, error)
	UpdateDevice(id uint64, config device.Config) error
	DeleteDevice(id uint64) error
	Datetime
	WiFiServer
	WiFiClient
	WAN
	DNS
	GPIO
}
