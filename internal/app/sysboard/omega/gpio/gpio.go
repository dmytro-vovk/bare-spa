//go:build sysboard

package gpio

import (
	"github.com/Sergii-Kirichok/pr/internal/app/types"
	"github.com/Sergii-Kirichok/pr/pkg/binary"
	"github.com/Sergii-Kirichok/pr/pkg/memory"
	log "github.com/sirupsen/logrus"
)

// todo: test with modes equal <nil> slice
func Setup(gpios map[byte]*types.GPIO, modes map[string]types.ModeBits) {
	const modeValue = 0b01

	mem := memory.At(baseAddr, baseSize)
	value := mem.ReadWord(modeOffset01)
	for desc, mod := range modes {
		log.Printf("Setting %s\n", desc)
		value = binary.SetFromTo(mod.FirstBit, mod.LastBit, value, modeValue)
	}
	mem.WriteWord(modeOffset01, value)
	mem.Close()

	for _, io := range gpios {
		SetDirection(io.Pin, io.Direction)
		SetLevel(io.Pin, io.Level)
		SetInversion(io.Pin, io.Inversion)
	}
}

func getParameterBit(pin byte, offset binary.Bits) bool {
	mem := memory.At(baseAddr, baseSize)
	defer mem.Close()

	value := mem.ReadWord(offset)

	mask := binary.GetBits(uint(pin))
	return binary.HasBits(value, mask)
}

func GetDirection(pin byte) types.DirectionBit {
	return types.DirectionBit(getParameterBit(getDirectionOffset(pin)))
}

func getDirectionOffset(pin byte) (byte, binary.Bits) {
	p, offset := pin, binary.Bits(0)
	switch {
	case pin <= 31:
		offset = ctrl0Offset
	case pin >= 32 && pin <= 63:
		offset = ctrl1Offset
		p -= 32
	case pin >= 64 && pin <= 95:
		offset = ctrl2Offset
		p -= 64
	}

	return p, offset
}

func SetDirection(pin byte, bit types.DirectionBit) {
	mem := memory.At(baseAddr, baseSize)
	defer mem.Close()

	pin, offset := getDirectionOffset(pin)
	value := mem.ReadWord(offset)
	mask := binary.GetBits(uint(pin))
	if bit {
		value = binary.SetBits(value, mask)
	} else {
		value = binary.ClearBits(value, mask)
	}
	mem.WriteWord(offset, value)
}

func GetInversion(pin byte) types.InversionBit {
	return types.InversionBit(getParameterBit(getInversionOffset(pin)))
}

func getInversionOffset(pin byte) (byte, binary.Bits) {
	p, offset := pin, binary.Bits(0)
	switch {
	case pin <= 31:
		offset = pol0Offset
	case pin >= 32 && pin <= 63:
		offset = pol1Offset
		p -= 32
	case pin >= 64 && pin <= 95:
		offset = pol2Offset
		p -= 64
	}

	return p, offset
}

// fixme: doesn't work
func SetInversion(pin byte, bit types.InversionBit) {
	mem := memory.At(baseAddr, baseSize)
	defer mem.Close()

	pin, offset := getInversionOffset(pin)
	value := mem.ReadWord(offset)
	mask := binary.GetBits(uint(pin))
	if bit {
		value = binary.SetBits(value, mask)
	} else {
		value = binary.ClearBits(value, mask)
	}
	mem.WriteWord(offset, value)
}

func GetLevel(pin byte) types.LevelBit {
	return types.LevelBit(getParameterBit(getLevelOffset(pin)))
}

func getLevelOffset(pin byte) (byte, binary.Bits) {
	p, offset := pin, binary.Bits(0)
	switch {
	case pin <= 31:
		offset = data0Offset
	case pin >= 32 && pin <= 63:
		offset = data1Offset
		p -= 32
	case pin >= 64 && pin <= 95:
		offset = data2Offset
		p -= 64
	}

	return p, offset
}

func SetLevel(pin byte, bit types.LevelBit) {
	mem := memory.At(baseAddr, baseSize)
	defer mem.Close()

	p, offset := pin, binary.Bits(0)
	switch {
	case pin <= 31:
		offset = dclr0Offset
		if bit {
			offset = dset0Offset
		}
	case pin >= 32 && pin <= 63:
		offset = dclr1Offset
		if bit {
			offset = dset1Offset
		}
		p -= 32
	case pin >= 64 && pin <= 95:
		offset = dclr2Offset
		if bit {
			offset = dset2Offset
		}
		p -= 64
	}

	mask := binary.GetBits(uint(p))
	mem.WriteWord(offset, mask)
}
