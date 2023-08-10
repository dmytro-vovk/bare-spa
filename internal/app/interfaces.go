package app

import (
	"fmt"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/i2c"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/spi"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/uart"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *Application) GetUARTConfig(req struct {
	ID byte `json:"id"`
}) (*uart.Config, error) {
	return a.board.GetUARTConfig(req.ID)
}

func (a *Application) SetUARTConfig(req struct {
	ID     byte        `json:"id"`
	Config uart.Config `json:"config"`
}) error {
	return a.notify(fmt.Sprintf("uart-%d.setConfig", req.ID), req.Config, a.board.SetUARTConfig(req.ID, req.Config))
}

func (a *Application) GetSPIConfig(req struct {
	ID byte `json:"id"`
}) (*spi.Config, error) {
	return a.board.GetSPIConfig(req.ID)
}

func (a *Application) SetSPIConfig(req struct {
	ID     byte       `json:"id"`
	Config spi.Config `json:"config"`
}) error {
	return a.notify(fmt.Sprintf("spi-%d.setConfig", req.ID), req.Config, a.board.SetSPIConfig(req.ID, req.Config))
}

func (a *Application) GetI2CConfig(req struct {
	ID byte `json:"id"`
}) (*i2c.Config, error) {
	return a.board.GetI2CConfig(req.ID)
}

func (a *Application) SetI2CConfig(req struct {
	ID     byte       `json:"id"`
	Config i2c.Config `json:"config"`
}) error {
	return a.notify(fmt.Sprintf("i2c-%d.setConfig", req.ID), req.Config, a.board.SetI2CConfig(req.ID, req.Config))
}

func (a *Application) GetTemperatures() (map[string]int16, error) {
	return a.board.GetTemperatures()
}

func (a *Application) observeTemperature() {
	log.Info("temperature observer was launched")
	for range time.NewTicker(time.Second).C {
		temp, err := a.board.GetTemperatures()
		if err != nil {
			log.Warningf("temperature observer error: %s", err)
			continue
		}

		a.Notify("system.setTemperatures", temp)
	}
}
