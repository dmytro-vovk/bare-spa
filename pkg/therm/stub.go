//go:build !sysboard

package therm

import (
	"math/rand"
)

func (t *Thermometer) Temperature(registerData string) (int16, error) {
	return int16(23 + rand.Intn(5) - rand.Intn(5)), nil
}
