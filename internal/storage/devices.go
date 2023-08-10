package storage

import "github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"

func (s *Storage) ListDevices() ([]*types.Device, error) {
	return s.devicesDB.All().([]*types.Device), nil
}

func (s *Storage) AddDevice(device *types.Device) error {
	id, err := s.devicesDB.Append(device)
	if err != nil {
		return err
	}
	device.ID = id
	return s.UpdateDevice(device)
}

func (s *Storage) UpdateDevice(device *types.Device) error {
	return s.devicesDB.Set(device.ID, device)
}

func (s *Storage) DeleteDevice(id int) error {
	return s.devicesDB.Delete(id)
}
