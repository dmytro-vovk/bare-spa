package therm

import (
	"time"
)

type Type string

const (
	Resistive    Type = "Резистивный"
	Infrared     Type = "Инфракрасный"
	Thermocouple Type = "Термопара на MLX90614"
	DS18B20      Type = "DS18B20"
	DTH          Type = "DTH-XX" // Температура и влажность
)

type Thermometer struct {
	Type             Type      `json:"type"` // Тип, на основе этого типа происходят расчёты показаний
	Name             string    `json:"name"`
	RegisterAddr     uint16    `json:"register_addr"` // Адрес регистра, по которому получаем данные
	RegisterQuantity uint16    `json:"register_quantity"`
	LastSeen         time.Time `json:"last_seen"` // Когда последний раз были получены данные
	Interval         int       `json:"interval"`  // Максимальный интервал между опросами этого датчика в секундах
	*Calibration     `json:"calibration"`
}

type Calibration struct {
	V0 int16 // Данные полученные при калибровке нуля
	Vc int16 // Данные полученные при втором значении температуры калибровки
	Tc int16 // Температура, при которой были получены данные второго значения калибровки
}

func (c *Calibration) SetZero(v int16) {
	c.V0 = v
}

func (c *Calibration) SetSecond(v, t int16) {
	c.Vc = v
	c.Tc = t
}
