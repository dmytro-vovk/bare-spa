package omega

import (
	"sort"

	"github.com/Sergii-Kirichok/pr/internal/app/device"
	"github.com/Sergii-Kirichok/pr/internal/app/device/address"
	"github.com/Sergii-Kirichok/pr/internal/app/device/module"
)

func (s *Sysboard) SupportedDevicesTypes() []device.Type {
	keys := make([]string, 0, len(s.SupportedDevices))
	for kind := range s.SupportedDevices {
		keys = append(keys, string(kind))
	}
	sort.Strings(keys)

	types := make([]device.Type, 0, len(keys))
	for _, kind := range keys {
		types = append(types, device.Type(kind))
	}

	return types
}

func (s *Sysboard) SupportedDeviceInterfaces(kind device.Type) []device.Interface {
	set := kind.Interfaces()
	keys := make([]string, 0, len(set))
	for iface := range kind.Interfaces() {
		keys = append(keys, string(iface))
	}
	sort.Strings(keys)

	ifaces := make([]device.Interface, 0, len(set))
	for _, iface := range keys {
		ifaces = append(ifaces, device.Interface(iface))
	}

	return ifaces
}

func (s *Sysboard) SupportedDeviceModules(kind device.Type) []module.Type {
	set := kind.Modules()
	keys := make([]string, 0, len(set))
	for mod := range kind.Modules() {
		keys = append(keys, string(mod))
	}
	sort.Strings(keys)

	modules := make([]module.Type, 0, len(set))
	for _, mod := range keys {
		modules = append(modules, module.Type(mod))
	}

	return modules
}

func (s *Sysboard) AvailableDevicesAddresses() []address.Address {
	return s.ConnectedDevices.AvailableAddresses()
}

func (s *Sysboard) ConnectedDevicesList(offset, limit int) ([]*device.Entity, error) {
	return s.ConnectedDevices.ReadMany(offset, limit)
}

func (s *Sysboard) ConnectedDevicesCount() int {
	return s.ConnectedDevices.Count()
}

func (s *Sysboard) CreateDevice(kind device.Type, config device.Config) (device.Device, error) {
	return s.ConnectedDevices.CreateOne(kind, config)
}

func (s *Sysboard) ReadDevice(id uint64) (device.Device, error) {
	return s.ConnectedDevices.ReadOne(id)
}

func (s *Sysboard) UpdateDevice(id uint64, config device.Config) error {
	return s.ConnectedDevices.UpdateOne(id, config)
}

func (s *Sysboard) DeleteDevice(id uint64) error {
	return s.ConnectedDevices.DeleteOne(id)
}
