package perimeth

import (
	"github.com/Sergii-Kirichok/pr/internal/app/device"
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard/omega"
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard/omega/gpio"
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard/omega/network"
	"github.com/Sergii-Kirichok/pr/internal/app/types"
	"github.com/pkg/errors"
)

type Sysboard struct {
	*omega.Sysboard
}

func New(name string) (*Sysboard, error) {
	sb, err := omega.New(name)
	if err != nil {
		return nil, err
	}

	sb.GPIO = map[byte]*types.GPIO{
		4: {
			Name:      "PIN-4",
			Pin:       4,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		5: {
			Name:      "PIN-5",
			Pin:       5,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		6: {
			Name:      "PIN-6",
			Pin:       6,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		8: {
			Name:      "PIN-8",
			Pin:       8,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		9: {
			Name:      "PIN-9",
			Pin:       9,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		11: {
			Name:      "PIN-11",
			Pin:       11,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		14: {
			Name:      "PIN-14",
			Pin:       14,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		15: {
			Name:      "PIN-15",
			Pin:       15,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		42: {
			Name:      "Relay",
			Pin:       42,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
	}

	sb.SupportedDevices = map[device.Type]struct{}{
		device.TypeIncubator1: {},
		device.TypeIncubator2: {},
		device.TypeEthIO2:     {},
		device.TypeEthIO4:     {},
		device.TypeUPS12:      {},
	}

	return &Sysboard{Sysboard: sb}, nil
}

func (s *Sysboard) Setup() error {
	gpio.Setup(s.GPIO, map[string]types.ModeBits{
		"GPIO mode": {
			FirstBit: 0,
			LastBit:  1,
		},
		"I2S GPIO mode": {
			FirstBit: 6,
			LastBit:  7,
		},
		"UART2 GPIO mode": {
			FirstBit: 26,
			LastBit:  27,
		},
		"PWM1 GPIO mode": {
			FirstBit: 30,
			LastBit:  31,
		},
	})

	err := network.Setup(s.Network)
	return errors.Wrap(err, "omega: unable to setup sysboard")
}
