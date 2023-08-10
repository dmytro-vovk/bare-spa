package device

import (
	"fmt"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device/address"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device/module"
)

type Type string

const (
	TypeEthIO2     Type = "Eth-IO-2"
	TypeEthIO4     Type = "Eth-IO-4"
	TypeIncubator1 Type = "Incubator-1"
	TypeIncubator2 Type = "Incubator-2"
	TypeOmega1     Type = "Omega-1"
	TypeUPS12      Type = "UPS12"
	TypeEmbedded   Type = "Embedded"
)

type Config struct {
	Name      string          `json:"name"`
	Interface Interface       `json:"interface"`
	Address   address.Address `json:"address"`
}

type Device interface {
	ID() uint64
	Type() Type
	Address() address.Address
	Entity() *Entity
	Update(config Config) error
	Delete() error
	CreateModule(kind module.Type, config module.Config) (module.Module, error)
	ReadModule(id uint64) (module.Module, error)
	UpdateModule(id uint64, config module.Config) error
	DeleteModule(id uint64) error
}

type ErrUnknownType struct {
	Kind Type
}

func (e *ErrUnknownType) Error() string {
	return fmt.Sprintf("unknown device type %q", e.Kind)
}

func newDevice(id *uint64, kind Type, config Config) (Device, error) {
	builder, ok := map[Type]func(*uint64, Config) Device{
		TypeEthIO2:     newEthIO2,
		TypeEthIO4:     newEthIO4,
		TypeIncubator1: newIncubator1,
		TypeIncubator2: newIncubator2,
		TypeOmega1:     newOmega1,
		TypeUPS12:      newUPS12,
		TypeEmbedded:   newEmbedded,
	}[kind]
	if !ok {
		return nil, &ErrUnknownType{Kind: kind}
	}

	if !kind.IsInterfaceSupported(config.Interface) {
		return nil, &ErrInterfaceNotSupported{
			DeviceType:      kind,
			DeviceInterface: config.Interface,
		}
	}

	if err := address.Borrow(config.Address); err != nil {
		return nil, err
	}

	return builder(id, config), nil
}

type device struct {
	id      uint64
	name    string
	kind    Type
	modules module.Modules
	iface   Interface
	addr    address.Address
}

func (d *device) ID() uint64 {
	return d.id
}

func (d *device) Type() Type {
	return d.kind
}

func (d *device) Address() address.Address {
	return d.addr
}

func (d *device) Entity() *Entity {
	return &Entity{
		ID:        d.id,
		Name:      d.name,
		Type:      d.kind,
		Interface: d.iface,
		Address:   d.addr,
		Modules:   d.modules.ReadAll(),
	}
}

func (d *device) Update(config Config) error {
	if d.Type() == TypeEmbedded {
		return &ErrEmbeddedDeviceType{ID: d.id}
	}

	if !d.kind.IsInterfaceSupported(config.Interface) {
		return &ErrInterfaceNotSupported{
			DeviceType:      d.Type(),
			DeviceInterface: config.Interface,
		}
	}

	if d.Address() != config.Address {
		if err := address.Borrow(config.Address); err != nil {
			return err
		}
		_ = address.Release(d.Address())
	}

	d.name = config.Name
	d.iface = config.Interface
	d.addr = config.Address

	return nil
}

func (d *device) Delete() error {
	if d.Type() == TypeEmbedded {
		return &ErrEmbeddedDeviceType{ID: d.id}
	}

	_ = address.Release(d.Address())

	return nil
}

func (d *device) CreateModule(kind module.Type, config module.Config) (module.Module, error) {
	if !d.kind.IsModuleSupported(kind) {
		return nil, &ErrModuleNotSupported{
			DeviceType: d.Type(),
			ModuleType: kind,
		}
	}

	return d.modules.CreateOne(kind, config)
}

func (d *device) ReadModule(id uint64) (module.Module, error) {
	return d.modules.ReadOne(id)
}

func (d *device) UpdateModule(id uint64, config module.Config) error {
	return d.modules.UpdateOne(id, config)
}

func (d *device) DeleteModule(id uint64) error {
	return d.modules.DeleteOne(id)
}
