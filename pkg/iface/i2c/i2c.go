package i2c

import "fmt"

type I2C struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Config Config `json:"config"`
}

type Config struct {
	Enabled    bool `json:"enabled"`
	BaudRate   int  `json:"baudRate"`
	Addressing int  `json:"addressing"`
}

var (
	baudRateValues = map[int]struct{}{
		10_000:    {},
		100_000:   {},
		400_000:   {},
		1_000_000: {},
		3_400_000: {},
	}

	addressingValues = map[int]struct{}{
		7:  {},
		16: {},
	}
)

func (i *I2C) SetConfig(config Config) error {
	if _, ok := baudRateValues[config.BaudRate]; !ok {
		return fmt.Errorf("invalid baud rate value '%d'", config.BaudRate)
	}

	if _, ok := addressingValues[config.Addressing]; !ok {
		return fmt.Errorf("invalid addressing value '%d'", config.Addressing)
	}

	i.Config = config
	return nil
}
