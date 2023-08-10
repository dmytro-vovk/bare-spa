package app

import (
	"github.com/Sergii-Kirichok/pr/internal/app/device"
	"github.com/Sergii-Kirichok/pr/internal/app/device/address"
	"github.com/Sergii-Kirichok/pr/internal/app/device/module"
	"github.com/Sergii-Kirichok/pr/internal/app/ifaces"
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard"
	"github.com/Sergii-Kirichok/pr/internal/app/types"
)

func (a *Application) GetSysboard() (ifaces.Sysboard, error) {
	return a.board, nil
}

func (a *Application) GetSystemBoardTypesList() ([]string, error) {
	return []string{
		sysboard.Omega,
		sysboard.PerimEth,
		sysboard.MiniCom,
	}, nil
}

// GetModules TODO after cleanup of structs need update this function
func (a *Application) GetModules() (map[types.DeviceType]types.Device, error) {
	return types.DeviceListBase, nil
}

func (a *Application) SupportedDevicesTypes() ([]device.Type, error) {
	return a.board.SupportedDevicesTypes(), nil
}

func (a *Application) SupportedDeviceInterfaces(req struct {
	Type device.Type `json:"type"`
}) ([]device.Interface, error) {
	return a.board.SupportedDeviceInterfaces(req.Type), nil
}

func (a *Application) SupportedDeviceModules(req struct {
	Type device.Type `json:"type"`
}) ([]module.Type, error) {
	return a.board.SupportedDeviceModules(req.Type), nil
}

func (a *Application) AvailableDevicesAddresses() ([]address.Address, error) {
	return a.board.AvailableDevicesAddresses(), nil
}

func (a *Application) ConnectedDevicesList(req struct {
	Page  int `json:"page" validate:"gte=1"`
	Limit int `json:"limit" validate:"oneof=1 2 3"`
}) ([]*device.Entity, error) {
	offset := (req.Page - 1) * req.Limit
	return a.board.ConnectedDevicesList(offset, req.Limit)
}

func (a *Application) ConnectedDevicesCount() (int, error) {
	return a.board.ConnectedDevicesCount(), nil
}

func (a *Application) CreateDevice(req struct {
	Type   device.Type   `json:"type"`
	Config device.Config `json:"config"`
}) error {
	_, err := a.board.CreateDevice(req.Type, req.Config)
	if err != nil {
		return err
	}

	return a.notify("devices.create", struct{}{}, nil)
}

func (a *Application) ReadDevice(req struct {
	ID uint64 `json:"id"`
}) (*device.Entity, error) {
	dev, err := a.board.ReadDevice(req.ID)
	if err != nil {
		return nil, err
	}

	return dev.Entity(), nil
}

func (a *Application) UpdateDevice(req struct {
	ID     uint64        `json:"id"`
	Config device.Config `json:"config"`
}) error {
	err := a.board.UpdateDevice(req.ID, req.Config)
	return a.notify("devices.update", struct{}{}, err)
}

func (a *Application) DeleteDevice(req struct {
	ID uint64 `json:"id"`
}) error {
	return a.notify("devices.delete", struct{}{}, a.board.DeleteDevice(req.ID))
}

func (a *Application) CreateModule(req struct {
	DeviceID uint64        `json:"deviceID"`
	Type     module.Type   `json:"type"`
	Config   module.Config `json:"config"`
}) error {
	dev, err := a.board.ReadDevice(req.DeviceID)
	if err != nil {
		return err
	}

	_, err = dev.CreateModule(req.Type, req.Config)
	if err != nil {
		return err
	}

	return a.notify("devices.createModule", struct{}{}, nil)
}

func (a *Application) ReadModule(req struct {
	DeviceID uint64 `json:"deviceID"`
	ModuleID uint64 `json:"moduleID"`
}) (*module.Entity, error) {
	dev, err := a.board.ReadDevice(req.DeviceID)
	if err != nil {
		return nil, err
	}

	mod, err := dev.ReadModule(req.ModuleID)
	if err != nil {
		return nil, err
	}

	return mod.Entity(), nil
}

func (a *Application) UpdateModule(req struct {
	DeviceID uint64        `json:"deviceID"`
	ModuleID uint64        `json:"moduleID"`
	Config   module.Config `json:"config"`
}) error {
	dev, err := a.board.ReadDevice(req.DeviceID)
	if err != nil {
		return err
	}

	err = dev.UpdateModule(req.ModuleID, req.Config)

	return a.notify("devices.updateModule", struct{}{}, err)
}

func (a *Application) DeleteModule(req struct {
	DeviceID uint64 `json:"deviceID"`
	ModuleID uint64 `json:"moduleID"`
}) error {
	dev, err := a.board.ReadDevice(req.DeviceID)
	if err != nil {
		return err
	}

	err = dev.DeleteModule(req.ModuleID)

	return a.notify("devices.deleteModule", struct{}{}, err)
}
