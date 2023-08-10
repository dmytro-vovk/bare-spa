package minicom

import (
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
		16: {
			Name:      "PIN-16",
			Pin:       16,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		17: {
			Name:      "PIN-17",
			Pin:       17,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		18: {
			Name:      "PIN-18",
			Pin:       18,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
		19: {
			Name:      "PIN-19",
			Pin:       19,
			Level:     types.LevelLow,
			Direction: types.DirectionOutput,
			Inversion: types.NonInverted,
		},
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
