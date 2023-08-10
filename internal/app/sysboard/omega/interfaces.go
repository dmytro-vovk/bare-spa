package omega

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/i2c"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/spi"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/uart"
)

func (s *Sysboard) GetInterfaces() []iface.Interface {
	interfaces := s.getUARTs()
	interfaces = append(interfaces, s.getSPIs()...)
	interfaces = append(interfaces, s.getI2Cs()...)
	return interfaces
}

func (s *Sysboard) GetUARTConfig(id byte) (*uart.Config, error) {
	i, err := s.getUART(id)
	if err != nil {
		return nil, err
	}

	return &i.Config, nil
}

func (s *Sysboard) SetUARTConfig(id byte, cfg uart.Config) error {
	i, err := s.getUART(id)
	if err != nil {
		return err
	}

	return errors.BadRequest.Use(i.SetConfig(cfg))
}

func (s *Sysboard) GetTemperatures() (map[string]int16, error) {
	const (
		masterID = 1 // UART number
		slaveID  = 1 // device ID
	)

	i, err := s.getUART(masterID)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int16)
	for _, sensor := range s.Thermometers {
		res, err := i.Connect(slaveID).ReadInputRegisters(sensor.RegisterAddr, sensor.RegisterQuantity)
		if err != nil {
			return nil, err
		}

		temp, err := sensor.Temperature(res)
		if err != nil {
			return nil, err
		}

		results[sensor.Name] = temp
	}

	return results, nil
}

func (s *Sysboard) getUARTs() []iface.Interface {
	interfaces := make([]iface.Interface, 0, len(s.UART))
	for n, i := range s.UART {
		interfaces = append(interfaces, iface.Interface{
			Name:   i.Name,
			Type:   iface.UART,
			Number: n,
		})
	}

	return interfaces
}

func (s *Sysboard) getUART(id byte) (*uart.UART, error) {
	i, ok := s.UART[id]
	if !ok {
		return nil, errors.BadRequest.Newf("UART-%d not found", id)
	}

	return i, nil
}

func (s *Sysboard) GetSPIConfig(id byte) (*spi.Config, error) {
	i, err := s.getSPI(id)
	if err != nil {
		return nil, err
	}

	return &i.Config, nil
}

func (s *Sysboard) SetSPIConfig(id byte, cfg spi.Config) error {
	i, err := s.getSPI(id)
	if err != nil {
		return err
	}

	return errors.BadRequest.Use(i.SetConfig(cfg))
}

func (s *Sysboard) getSPIs() []iface.Interface {
	interfaces := make([]iface.Interface, 0, len(s.SPI))
	for n, i := range s.SPI {
		interfaces = append(interfaces, iface.Interface{
			Name:   i.Name,
			Type:   iface.SPI,
			Number: n,
		})
	}

	return interfaces
}

func (s *Sysboard) getSPI(id byte) (*spi.SPI, error) {
	i, ok := s.SPI[id]
	if !ok {
		return nil, errors.BadRequest.Newf("SPI-%d not found", id)
	}

	return i, nil
}

func (s *Sysboard) GetI2CConfig(id byte) (*i2c.Config, error) {
	i, err := s.getI2C(id)
	if err != nil {
		return nil, err
	}

	return &i.Config, nil
}

func (s *Sysboard) SetI2CConfig(id byte, cfg i2c.Config) error {
	i, err := s.getI2C(id)
	if err != nil {
		return err
	}

	return errors.BadRequest.Use(i.SetConfig(cfg))
}

func (s *Sysboard) getI2Cs() []iface.Interface {
	interfaces := make([]iface.Interface, 0, len(s.I2C))
	for n, i := range s.I2C {
		interfaces = append(interfaces, iface.Interface{
			Name:   i.Name,
			Type:   iface.I2C,
			Number: n,
		})
	}

	return interfaces
}

func (s *Sysboard) getI2C(id byte) (*i2c.I2C, error) {
	i, ok := s.I2C[id]
	if !ok {
		return nil, errors.BadRequest.Newf("I2C-%d not found", id)
	}

	return i, nil
}
