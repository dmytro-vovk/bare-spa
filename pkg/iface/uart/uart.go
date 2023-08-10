package uart

import (
	"fmt"
)

type UART struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Config Config `json:"config"`
}

type Config struct {
	Enabled  bool   `json:"enabled"`
	BaudRate int    `json:"baudRate"`
	DataBits int    `json:"dataBits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stopBits"`
}

var (
	baudRateValues = map[int]struct{}{
		9_600:   {},
		14_400:  {},
		19_200:  {},
		31_250:  {},
		38_400:  {},
		56_000:  {},
		57_600:  {},
		76_800:  {},
		115_200: {},
		128_000: {},
		230_400: {},
		250_000: {},
		256_000: {},
	}

	dataBitsValues = map[int]struct{}{
		5: {},
		6: {},
		7: {},
		8: {},
	}

	parityValues = map[string]struct{}{
		"E": {},
		"M": {},
		"N": {},
		"O": {},
		"S": {},
	}

	stopBitsValues = map[int]struct{}{
		1: {},
		2: {},
	}
)

func (u *UART) SetConfig(config Config) error {
	if _, ok := baudRateValues[config.BaudRate]; !ok {
		return fmt.Errorf("invalid baud rate value '%d'", config.BaudRate)
	}

	if _, ok := dataBitsValues[config.DataBits]; !ok {
		return fmt.Errorf("invalid data bits value '%d'", config.DataBits)
	}

	if _, ok := parityValues[config.Parity]; !ok {
		return fmt.Errorf("invalid parity value '%s'", config.Parity)
	}

	if _, ok := stopBitsValues[config.StopBits]; !ok {
		return fmt.Errorf("invalid stop bits value '%d'", config.StopBits)
	}

	u.Config = config
	return nil
}
