package iface

type Type string

const (
	UART  Type = "UART"
	SPI   Type = "SPI"
	I2C   Type = "I2C"
	RS485 Type = "RS485"
	Radio Type = "RADIO"
)

type Interface struct {
	Name   string `json:"name"`
	Type   Type   `json:"type"`
	Number byte   `json:"number"`
}
