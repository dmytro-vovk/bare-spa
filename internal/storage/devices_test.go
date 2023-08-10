package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO need update tests when API will be stable
func TestDevices(t *testing.T) {
	db := "test.db"
	_ = os.Remove(db)
	s := New(db)
	if list, err := s.ListDevices(); assert.NoError(t, err) {
		assert.Empty(t, list)
	}
	d := PackedDevice{
		Name: "Name",
		Type: "Type",
		Data: nil,
	}
	if assert.NoError(t, s.AddNewDevice(&d)) {
		if list, err := s.ListDevices(); assert.NoError(t, err) {
			assert.Equal(t, 0, len(list))
			assert.Equal(t, PackedDevice{ID: 0, Name: "Name", Type: "Type", Data: json.RawMessage(nil)}, d)
		}
		d.Name = "Name 2"
		d.Type = "Type 2"
		// d.Interface = "IF 2"
		// d.DeviceID = 2
		// if assert.NoError(t, s.UpdateDevice(&d)) {
		// 	if devs, err := s.ConnectedDevicesList(); assert.NoError(t, err) {
		// 		assert.Equal(t, &types.Device{
		// 			ID:        1,
		// 			Name:      "Name 2",
		// 			Type:      "Type 2",
		// 			Interface: "IF 2",
		// 			DeviceID:  2,
		// 		}, devs[0])
		// 	}
		// }
		assert.NoError(t, s.DeleteDevice(1))
		if list, err := s.ListDevices(); assert.NoError(t, err) {
			assert.Empty(t, list)
		}
	}
	_ = os.RemoveAll(db)
}
