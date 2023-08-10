package sysboard

import (
	"fmt"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/ifaces"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/minicom"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/omega"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/perimeth"
)

const (
	Omega    = "Omega"
	PerimEth = "PerimEth"
	MiniCom  = "MiniCom"
)

func New(boardType string) (ifaces.Sysboard, error) {
	const errPrefix = "system board:"

	switch boardType {
	case Omega:
		sysboard, err := omega.New(Omega)
		if err != nil {
			return nil, fmt.Errorf("%s %w", errPrefix, err)
		}

		return load(sysboard), err
	case PerimEth:
		sysboard, err := perimeth.New(PerimEth)
		if err != nil {
			return nil, fmt.Errorf("%s %w", errPrefix, err)
		}

		return load(sysboard), err
	case MiniCom:
		sysboard, err := minicom.New(MiniCom)
		if err != nil {
			return nil, fmt.Errorf("%s %w", errPrefix, err)
		}

		return load(sysboard), err
	default:
		return nil, errors.ParsingErr.Newf("unknown system board type %s", boardType)
	}
}

func Reset(s ifaces.Sysboard) error { return deleteConfig(s) }
