package storage

import (
	"encoding/json"
)

type PackedDevice struct {
	ID   int             `json:"id"`
	Name string          `json:"name"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (s *Storage) ListNewDevices() ([]*PackedDevice, error) {
	return s.newDevicesDB.All().([]*PackedDevice), nil
}

func (s *Storage) AddNewDevice(device *PackedDevice) error {
	_, err := s.newDevicesDB.Append(device)
	return err
}
