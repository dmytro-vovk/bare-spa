package omega

import (
	"fmt"
	"github.com/Sergii-Kirichok/pr/internal/app/errors"
	"github.com/Sergii-Kirichok/pr/internal/app/sysboard/omega/gpio"
	"github.com/Sergii-Kirichok/pr/internal/app/types"
	log "github.com/sirupsen/logrus"
	"sort"
)

type ErrGPIONotUsed struct {
	Pin byte
}

func (e ErrGPIONotUsed) Error() string {
	return fmt.Sprintf("GPIO[%2d] not used", e.Pin)
}

func (s *Sysboard) GetGPIO() []types.GPIO {
	pins := make([]int, 0, len(s.GPIO))
	for pin := range s.GPIO {
		pins = append(pins, int(pin))
	}

	sort.Ints(pins)
	gpios := make([]types.GPIO, 0, len(s.GPIO))
	for _, pin := range pins {
		gpios = append(gpios, *s.GPIO[byte(pin)])
	}

	return gpios
}

func (s *Sysboard) GetGPIODirection(pin byte) (types.DirectionBit, bool, error) {
	io, ok := s.GPIO[pin]
	if !ok {
		return false, false, errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to get GPIO direction")
	}

	var isInnerUpdate bool
	if dir := gpio.GetDirection(pin); io.Direction != dir {
		io.Direction = dir
		isInnerUpdate = true
		log.Printf("GPIO[%2d] direction was updated to %q", pin, dir)
	}

	return io.Direction, isInnerUpdate, nil
}

func (s *Sysboard) SetGPIODirection(pin byte, bit types.DirectionBit) error {
	io, ok := s.GPIO[pin]
	if !ok {
		return errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to set GPIO direction")
	}

	gpio.SetDirection(pin, bit)
	io.Direction = bit
	log.Printf("GPIO[%2d] direction was set as %q", pin, bit)
	return nil
}

func (s *Sysboard) GetGPIOInversion(pin byte) (types.InversionBit, bool, error) {
	io, ok := s.GPIO[pin]
	if !ok {
		return false, false, errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to get GPIO inversion")
	}

	var isInnerUpdate bool
	if inv := gpio.GetInversion(pin); io.Inversion != inv {
		io.Inversion = inv
		isInnerUpdate = true
		log.Printf("GPIO[%2d] inversion was updated to %q", pin, inv)
	}

	return io.Inversion, isInnerUpdate, nil
}

func (s *Sysboard) SetGPIOInversion(pin byte, bit types.InversionBit) error {
	io, ok := s.GPIO[pin]
	if !ok {
		return errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to set GPIO inversion")
	}

	gpio.SetInversion(pin, bit)
	io.Inversion = bit
	log.Printf("GPIO[%2d] inversion was set as %q", pin, bit)
	return nil
}

func (s *Sysboard) GetGPIOLevel(pin byte) (types.LevelBit, bool, error) {
	io, ok := s.GPIO[pin]
	if !ok {
		return false, false, errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to get GPIO level")
	}

	var isInnerUpdate bool
	if lvl := gpio.GetLevel(pin); io.Level != lvl {
		io.Level = lvl
		isInnerUpdate = true
		log.Printf("GPIO[%2d] level was updated to %q", pin, lvl)
	}

	return io.Level, isInnerUpdate, nil
}

func (s *Sysboard) SetGPIOLevel(pin byte, bit types.LevelBit) error {
	io, ok := s.GPIO[pin]
	if !ok {
		return errors.BadRequest.Wrap(ErrGPIONotUsed{Pin: pin}, "failed to set GPIO level")
	}

	gpio.SetLevel(pin, bit)
	io.Level = bit
	log.Printf("GPIO[%2d] level was set as %q", pin, bit)
	return nil
}
