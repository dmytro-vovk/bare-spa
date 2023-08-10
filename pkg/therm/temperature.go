//go:build sysboard

package therm

import (
	"github.com/pkg/errors"
	"strconv"
)

func (t *Thermometer) Temperature(registerData string) (int16, error) {
	v, err := strconv.ParseInt(registerData, 16, 16)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to get temperature from data 0x%X", v)
	}

	return t.temperature(int16(v)), nil
}

func (c *Calibration) temperature(v int16) int16 {
	cost := (c.V0 - c.Vc) / c.Tc
	return (c.V0 - v) / cost
}
