//go:build !sysboard

package gpio

import "github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"

var bits = map[byte]*struct {
	types.DirectionBit
	types.InversionBit
	types.LevelBit
}{}

func Setup(map[byte]*types.GPIO, map[string]types.ModeBits) {}

func GetDirection(pin byte) types.DirectionBit {
	if io := bits[pin]; io != nil {
		return io.DirectionBit
	}

	initGPIO(pin)
	return GetDirection(pin)
}

func SetDirection(pin byte, bit types.DirectionBit) { bits[pin].DirectionBit = bit }

func GetInversion(pin byte) types.InversionBit {
	if io := bits[pin]; io != nil {
		return io.InversionBit
	}

	initGPIO(pin)
	return GetInversion(pin)
}

func SetInversion(pin byte, bit types.InversionBit) { bits[pin].InversionBit = bit }

func GetLevel(pin byte) types.LevelBit {
	if io := bits[pin]; io != nil {
		return io.LevelBit
	}

	initGPIO(pin)
	return GetLevel(pin)
}

func SetLevel(pin byte, bit types.LevelBit) { bits[pin].LevelBit = bit }

func initGPIO(pin byte) {
	bits[pin] = &struct {
		types.DirectionBit
		types.InversionBit
		types.LevelBit
	}{
		DirectionBit: types.DirectionOutput,
		InversionBit: types.NonInverted,
		LevelBit:     types.LevelLow,
	}
}
