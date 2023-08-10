package device

import (
	"fmt"
	"sync"

	"github.com/Sergii-Kirichok/pr/internal/app/device/address"
	"github.com/Sergii-Kirichok/pr/internal/app/device/module"
	"github.com/Sergii-Kirichok/pr/internal/app/errors"
)

type Devices struct {
	id      uint64
	devices []Device
	mu      sync.Mutex
}

func NewDevices(devices ...Device) *Devices {
	return &Devices{
		devices: devices,
	}
}

func (d *Devices) CreateOne(kind Type, config Config) (Device, error) {
	dev, err := newDevice(&d.id, kind, config)
	if err != nil {
		return nil, errors.Wrap(err, "create device")
	}

	return d.add(dev), nil
}

func (d *Devices) add(dev Device) Device {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.devices = append(d.devices, dev)
	return dev
}

type ErrDeviceNotFound struct {
	ID uint64
}

func (e *ErrDeviceNotFound) Error() string {
	return fmt.Sprintf("device with id %d not found", e.ID)
}

func (d *Devices) ReadOne(id uint64) (Device, error) {
	idx, dev := d.find(id)
	if idx == -1 {
		return nil, errors.Wrap(&ErrDeviceNotFound{ID: id}, "read device")
	}

	return dev, nil
}

func (d *Devices) find(id uint64) (int, Device) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for idx, dev := range d.devices {
		if id == dev.ID() {
			return idx, dev
		}
	}

	return -1, nil
}

type ErrEmbeddedDeviceType struct {
	ID uint64
}

func (e *ErrEmbeddedDeviceType) Error() string {
	return fmt.Sprintf("device with id %d has embedded type", e.ID)
}

func (d *Devices) UpdateOne(id uint64, config Config) error {
	const errFormat = "update device"

	idx, dev := d.find(id)
	if idx == -1 {
		return errors.Wrap(&ErrDeviceNotFound{ID: id}, errFormat)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	err := dev.Update(config)
	return errors.Wrap(err, errFormat)
}

func (d *Devices) DeleteOne(id uint64) error {
	const errFormat = "delete device"

	idx, dev := d.find(id)
	if idx == -1 {
		return errors.Wrap(&ErrDeviceNotFound{ID: id}, errFormat)
	}

	if err := dev.Delete(); err != nil {
		return errors.Wrap(err, errFormat)
	}

	d.delete(idx)
	return nil
}

func (d *Devices) delete(index int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.devices = append(d.devices[:index], d.devices[index+1:]...)
}

type Entity struct {
	ID        uint64           `json:"id"`
	Name      string           `json:"name"`
	Type      Type             `json:"type"`
	Interface Interface        `json:"interface"`
	Address   address.Address  `json:"address"`
	Modules   []*module.Entity `json:"modules"`
}

func (d *Devices) ReadMany(offset, limit int) ([]*Entity, error) {
	const errPrefix = "devices list:"

	if offset < 0 {
		return nil, fmt.Errorf("%s offset can't be less than 0", errPrefix)
	} else if offset != 0 && offset >= d.Count() {
		return nil, fmt.Errorf("%s offset=%d greater than devices count", errPrefix, offset)
	}

	if limit < 0 {
		return nil, fmt.Errorf("%s limit can't be less than 0", errPrefix)
	} else if count := d.Count(); offset+limit > count {
		limit = count - offset
	}

	return d.list(offset, limit), nil
}

func (d *Devices) list(offset, limit int) []*Entity {
	d.mu.Lock()
	defer d.mu.Unlock()

	list := make([]*Entity, 0, limit)
	for _, dev := range d.devices[offset : offset+limit] {
		list = append(list, dev.Entity())
	}

	return list
}

func (d *Devices) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.devices)
}

func (d *Devices) AvailableAddresses() []address.Address {
	return address.ReadAll()
}
