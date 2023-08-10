package spi

import "fmt"

type SPI struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Config Config `json:"config"`
}

type Config struct {
	Enabled       bool   `json:"enabled"`
	BaudRate      int    `json:"baudRate"`
	ClockPolarity string `json:"clockPolarity"`
	ClockPhase    string `json:"clockPhase"`
	DataBits      int    `json:"dataBits"`
	FirstBit      string `json:"firstBit"`
}

var (
	baudRateValues = map[int]struct{}{
		31_250:    {},
		62_500:    {},
		128_000:   {},
		256_000:   {},
		512_000:   {},
		2_000_000: {},
		4_000_000: {},
	}

	clockPolarityValues = map[string]struct{}{
		"Low":  {},
		"High": {},
	}

	clockPhaseValues = map[string]struct{}{
		"1 Edge": {},
		"2 Edge": {},
	}

	dataBitsValues = map[int]struct{}{
		8:  {},
		16: {},
	}

	firstBitValues = map[string]struct{}{
		"MSB": {},
		"LSB": {},
	}
)

func (s *SPI) SetConfig(config Config) error {
	if _, ok := baudRateValues[config.BaudRate]; !ok {
		return fmt.Errorf("invalid baud rate value '%d'", config.BaudRate)
	}

	if _, ok := clockPolarityValues[config.ClockPolarity]; !ok {
		return fmt.Errorf("invalid clock polarity value '%s'", config.ClockPolarity)
	}

	if _, ok := clockPhaseValues[config.ClockPhase]; !ok {
		return fmt.Errorf("invalid clock phase value '%s'", config.ClockPhase)
	}

	if _, ok := dataBitsValues[config.DataBits]; !ok {
		return fmt.Errorf("invalid data bits value '%d'", config.DataBits)
	}

	if _, ok := firstBitValues[config.FirstBit]; !ok {
		return fmt.Errorf("invalid first bit value '%d'", config.DataBits)
	}

	s.Config = config
	return nil
}
