package storage

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"os"
	"strings"

	fdb "github.com/Sergii-Kirichok/DTekSpeachParser/internal/storage/file-db"
)

type Storage struct {
	path      string
	devicesDB *fdb.FileDB
	aliasesDB *fdb.FileDB

	newDevicesDB *fdb.FileDB // for devices in internal/app/devices
}

func New(path string) *Storage {
	return (&Storage{
		path: strings.TrimRight(path, "/") + "/",
	}).mustInit()
}

func (s *Storage) mustInit() *Storage {
	if err := os.MkdirAll(s.path, 0750); err != nil {
		panic(err)
	}
	s.devicesDB = fdb.Open(s.path+"devices.json", &types.Device{})
	s.aliasesDB = fdb.Open(s.path+"aliases.json", &types.Alias{})
	s.newDevicesDB = fdb.Open(s.path+"new-devices.json", &PackedDevice{})
	return s
}
