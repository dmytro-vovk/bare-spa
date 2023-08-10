package ifaces

import "github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"

type GPIO interface {
	GetGPIO() []types.GPIO
	GetGPIODirection(pin byte) (types.DirectionBit, bool, error)
	SetGPIODirection(pin byte, bit types.DirectionBit) error
	GetGPIOInversion(pin byte) (types.InversionBit, bool, error)
	SetGPIOInversion(pin byte, bit types.InversionBit) error
	GetGPIOLevel(pin byte) (types.LevelBit, bool, error)
	SetGPIOLevel(pin byte, bit types.LevelBit) error
}
