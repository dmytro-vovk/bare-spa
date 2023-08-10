package types

type Interface string

const (
	UART  Interface = "UART"
	SPI   Interface = "SPI"
	I2C   Interface = "I2C"
	RS485 Interface = "RS485"
	Radio Interface = "RADIO"
)

type UARTMode string

const (
	UARTModeNormal UARTMode = "Normal"
	UARTModeRS485  UARTMode = "RS485"
)

type SPIMode int

const (
	SPIMode0 SPIMode = iota
	SPIMode1
	SPIMode2
	SPIMode3
)

// COM "uart_1", "uart_2", "spi", "ethernet", "i2c", "Console"
type COM struct {
	Name           string    `json:"name"`           // Имя интерфейса
	Type           Interface `json:"type"`           // UART, SPI, I2C
	UartMode       UARTMode  `json:"uartMode"`       // UartMode || RS485Mode. Используем для определения текущего режима
	IsUARTAllowed  bool      `json:"isUARTAllowed"`  // Может ли порт работать в режиме UART
	IsRS485Allowed bool      `json:"isRS485Allowed"` // Может ли порт работать в режиме RS485
	SPIMode        SPIMode   `json:"spiMode"`        // Режим работы SPI
	Disabled       bool      `json:"disabled"`
	BaudRate       int       `json:"baudRate"`
	DataBits       int       `json:"dataBits"`
	Parity         byte      `json:"parity"`
	StopBits       byte      `json:"stopBits"`
	GPIO           uint8     `json:"gpio"` // Указывает либо chip select(cs), либо направление прием/передача
	Path           string    `json:"path"` // Путь к файлу устройства в разрезе COM порта и SPI
}
