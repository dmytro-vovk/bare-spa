package app

import (
	"github.com/Sergii-Kirichok/pr/internal/app/types"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *Application) observeGPIO() {
	log.Println("GPIO observer was launched")
	gpios := a.board.GetGPIO()
	for range time.NewTicker(time.Second).C {
		for _, io := range gpios {
			if dir, upd, _ := a.board.GetGPIODirection(io.Pin); upd {
				a.Notify("gpio.setDirection", struct {
					Pin       byte               `json:"pin"`
					Direction types.DirectionBit `json:"direction"`
				}{
					Pin:       io.Pin,
					Direction: dir,
				})
			}

			if inv, upd, _ := a.board.GetGPIOInversion(io.Pin); upd {
				a.Notify("gpio.setInversion", struct {
					Pin       byte               `json:"pin"`
					Inversion types.InversionBit `json:"inversion"`
				}{
					Pin:       io.Pin,
					Inversion: inv,
				})
			}

			if lvl, upd, _ := a.board.GetGPIOLevel(io.Pin); upd {
				a.Notify("gpio.setLevel", struct {
					Pin   byte           `json:"pin"`
					Level types.LevelBit `json:"level"`
				}{
					Pin:   io.Pin,
					Level: lvl,
				})
			}
		}
	}
}

func (a *Application) GetGPIO() ([]types.GPIO, error) {
	return a.board.GetGPIO(), nil
}

func (a *Application) SetGPIODirection(req struct {
	Pin       byte               `json:"pin"`
	Direction types.DirectionBit `json:"direction"`
}) error {
	return a.notify("gpio.setDirection", req, a.board.SetGPIODirection(req.Pin, req.Direction))
}

func (a *Application) SetGPIOInversion(req struct {
	Pin       byte               `json:"pin"`
	Inversion types.InversionBit `json:"inversion"`
}) error {
	return a.notify("gpio.setInversion", req, a.board.SetGPIOInversion(req.Pin, req.Inversion))
}

func (a *Application) SetGPIOLevel(req struct {
	Pin   byte           `json:"pin"`
	Level types.LevelBit `json:"level"`
}) error {
	return a.notify("gpio.setLevel", req, a.board.SetGPIOLevel(req.Pin, req.Level))
}
