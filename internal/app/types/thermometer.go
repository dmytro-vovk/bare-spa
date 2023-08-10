package types

import "time"

type ThermometerType string

const (
	Resistive    ThermometerType = "Резистивный"
	Infrared     ThermometerType = "Инфракрасный"
	Thermocouple ThermometerType = "Термопара на MLX90614"
	DS18B20      ThermometerType = "DS18B20"
	DTH          ThermometerType = "DTH-XX" // Температура и влажность
)

type ThermometerSettings struct {
	Name             string          `json:"name"`             // Имя
	Type             ThermometerType `json:"type"`             // Тип, на основе этого типа происходят расчёты показаний
	RegisterAddr     uint16          `json:"registerAddr"`     // Адрес регистра, по которому получаем данные
	RegisterData     uint16          `json:"registerData"`     // Фактические данные полученные из регистра
	Calibration0Val  uint16          `json:"calibration0Val"`  // Данные полученные при калибровке нуля
	Calibration1Val  uint16          `json:"calibration1Val"`  // Данные полученные при значении "Calibration1Temp"
	Calibration1Temp uint16          `json:"calibration1Temp"` // Температура, при которой были получены значения "Calibration1Val"
	LastSeen         time.Time       `json:"lastSeen"`         // Когда последний раз были получены данные
	Interval         int             `json:"interval"`         // Максимальный интервал между опросами этого датчика в секундах
}
