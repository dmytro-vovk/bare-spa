package omega

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/device/address"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/ifaces"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/omega/gpio"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/sysboard/omega/network"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/i2c"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/spi"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/iface/uart"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/therm"
)

var _ ifaces.Sysboard = (*Sysboard)(nil)

type Sysboard struct {
	Kind             string                      `json:"type"`    // Тип главной платы (Omega, PerimEth, MiniCom)
	Network          types.Network               `json:"network"` // Настройки сети и WiFi
	GPIO             map[byte]*types.GPIO        `json:"gpio"`    // Доступные GPIO на плате (Номера портов, доступных пользователю)
	UART             map[byte]*uart.UART         `json:"uart"`
	SPI              map[byte]*spi.SPI           `json:"spi"`
	I2C              map[byte]*i2c.I2C           `json:"i2c"`
	Thermometers     map[byte]*therm.Thermometer `json:"thermometers"`
	SupportedDevices map[device.Type]struct{}    `json:"supportedDevices"` // Список поддерживаемых устройств
	ConnectedDevices *device.Devices             `json:"connectedDevices"` // Список подключенных устройств
	//Modules          []*types.Module      `json:"modules"`          // Подключаемые модули (радио модуль, дополнительный порт RS485 или 3G модем)
	Datetime types.Datetime `json:"datetime"`
}

func (s *Sysboard) UnmarshalJSON(data []byte) error {
	var config Sysboard
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Network configuring
	if err := s.SetWAN(config.Network.WAN); err != nil {
		log.Println(err)
	}

	// GPIOs configuring
	for _, io := range config.GPIO {
		if err := s.SetGPIODirection(io.Pin, io.Direction); err != nil {
			log.Println(err)
		}

		if err := s.SetGPIOInversion(io.Pin, io.Inversion); err != nil {
			log.Println(err)
		}

		if err := s.SetGPIOLevel(io.Pin, io.Level); err != nil {
			log.Println(err)
		}
	}

	// Datetime configuring
	if err := s.SetDatetime(config.Datetime.Timestamp); err != nil {
		log.Println(err)
	}

	if err := s.SetNTP(config.Datetime.NTP); err != nil {
		log.Println(err)
	}

	return nil
}

func New(kind string) (*Sysboard, error) {
	const errFormat = "unable to create %s: %w"

	devices := device.NewDevices()

	_, err := devices.CreateOne(device.TypeEmbedded, device.Config{
		Name:      "Микроконтроллер",
		Interface: device.InterfaceUART,
		Address:   address.LeftBorder,
	})
	if err != nil {
		return nil, fmt.Errorf(errFormat, kind, err)
	}

	return &Sysboard{
		Kind: kind,
		Network: types.Network{
			WiFiAp: types.WiFiAp{
				Enabled:    false,
				SSID:       "Simple Solutions",
				Password:   "12345678",
				Channel:    "auto",
				TTL:        5,
				Encryption: types.WPA2,
				IP:         "192.168.1.1",
				Mask:       "255.255.255.0",
				DHCP:       types.DHCPOff,
			},
			WiFiCl: types.WiFiCl{
				Enabled: false,
				DHCP:    types.DHCPOn,
			},
			WAN: types.WAN{
				IP:      "192.168.0.10",
				Mask:    "255.255.255.0",
				Gateway: "192.168.0.1",
				DHCP:    types.DHCPOn,
			},
			DNS: types.DNS{
				DNS1: "192.168.0.1",
			},
		},
		GPIO: map[byte]*types.GPIO{
			0: {
				Name:      "PIN-0",
				Pin:       0,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			1: {
				Name:      "PIN-1",
				Pin:       1,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			2: {
				Name:      "PIN-2",
				Pin:       2,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			3: {
				Name:      "PIN-3",
				Pin:       3,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			11: {
				Name:      "PIN-11",
				Pin:       11,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			15: {
				Name:      "PIN-15",
				Pin:       15,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			16: {
				Name:      "PIN-16",
				Pin:       16,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
			17: {
				Name:      "PIN-17",
				Pin:       17,
				Level:     types.LevelLow,
				Direction: types.DirectionOutput,
				Inversion: types.NonInverted,
			},
		},
		UART: map[byte]*uart.UART{
			0: {
				Name: "Console",
				Path: "/dev/ttyS0",
				Config: uart.Config{
					BaudRate: 115_200,
					DataBits: 8,
					Parity:   "N",
					StopBits: 1,
				},
			},
			1: {
				Name: "UART-1",
				Path: "/dev/ttyS1",
				Config: uart.Config{
					BaudRate: 9_600,
					DataBits: 8,
					Parity:   "N",
					StopBits: 1,
				},
			},
		},
		SPI: map[byte]*spi.SPI{
			0: {
				Name: "SPI",
				Path: "/dev/spidev0.1",
				Config: spi.Config{
					BaudRate:      128_000,
					ClockPolarity: "Low",
					ClockPhase:    "1 Edge",
					DataBits:      8,
					FirstBit:      "MSB",
				},
			},
		},
		I2C: map[byte]*i2c.I2C{
			0: {
				Name: "I2C",
				Path: "/dev/i2c-0",
				Config: i2c.Config{
					BaudRate:   100_000,
					Addressing: 7,
				},
			},
		},
		Thermometers: map[byte]*therm.Thermometer{
			0: {
				Type:             therm.Resistive,
				Name:             "Резистивный #1",
				RegisterAddr:     0x0001,
				RegisterQuantity: 1,
				Interval:         5,
				Calibration: &therm.Calibration{
					V0: 0x0A40,
					Vc: 0x07E8,
					Tc: 25,
				},
			},
			1: {
				Type:             therm.Resistive,
				Name:             "Резистивный #2",
				RegisterAddr:     0x0002,
				RegisterQuantity: 1,
				Interval:         5,
				Calibration: &therm.Calibration{
					V0: 0x0A40,
					Vc: 0x07E8,
					Tc: 25,
				},
			},
			2: {
				Type:             therm.Infrared,
				Name:             "Инфракрасный",
				RegisterAddr:     0x0003,
				RegisterQuantity: 1,
				Interval:         5,
				Calibration: &therm.Calibration{
					V0: 0x0A40,
					Vc: 0x07E8,
					Tc: 25,
				},
			},
			3: {
				Type:             therm.Thermocouple,
				Name:             "Термопара-0",
				RegisterAddr:     0x0004,
				RegisterQuantity: 1,
				Interval:         5,
				Calibration: &therm.Calibration{
					V0: 0x0A40,
					Vc: 0x07E8,
					Tc: 25,
				},
			},
			4: {
				Type:             therm.Thermocouple,
				Name:             "Термопара-1",
				RegisterAddr:     0x0005,
				RegisterQuantity: 1,
				Interval:         5,
				Calibration: &therm.Calibration{
					V0: 0x0A40,
					Vc: 0x07E8,
					Tc: 25,
				},
			},
		},
		SupportedDevices: map[device.Type]struct{}{
			device.TypeOmega1: {},
			device.TypeEthIO2: {},
			device.TypeEthIO4: {},
		},
		ConnectedDevices: devices,
	}, nil
}

func (s *Sysboard) Type() string {
	return s.Kind
}

func (s *Sysboard) Setup() error {
	gpio.Setup(s.GPIO, map[string]types.ModeBits{
		"GPIO mode": {
			FirstBit: 0,
			LastBit:  1,
		},
		"I2S GPIO mode": {
			FirstBit: 6,
			LastBit:  7,
		},
		"UART2 GPIO mode": {
			FirstBit: 26,
			LastBit:  27,
		},
		"PWM1 GPIO mode": {
			FirstBit: 30,
			LastBit:  31,
		},
	})

	err := network.Setup(s.Network)
	return errors.Wrap(err, "omega: unable to setup sysboard")
}
