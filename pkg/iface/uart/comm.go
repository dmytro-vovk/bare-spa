//go:build sysboard

package uart

import (
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/goburrow/modbus"
)

type Conn struct {
	client modbus.Client
}

func (u *UART) Connect(id byte) *Conn {
	h := modbus.NewASCIIClientHandler(u.Path)
	h.BaudRate = u.Config.BaudRate
	h.DataBits = u.Config.DataBits
	h.Parity = u.Config.Parity
	h.StopBits = u.Config.StopBits
	h.SlaveId = id
	h.Timeout = time.Second

	return &Conn{client: modbus.NewClient(h)}
}

func (c *Conn) ReadInputRegisters(address, quantity uint16) (string, error) {
	res, err := c.client.ReadInputRegisters(address, quantity)
	if err != nil {
		return "", fmt.Errorf("error reading %d input register(s) at 0x%X: %w", quantity, address, err)
	}
	log.Infof("reading %d input register(s) at 0x%X result: 0x%X", quantity, address, res)

	return hex.EncodeToString(res), nil
}

func (c *Conn) ReadHoldingRegisters(address, quantity uint16) (results string, err error) {
	res, err := c.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return "", fmt.Errorf("error reading %d holding register(s) at 0x%X: %w", quantity, address, err)
	}
	log.Infof("reading %d holding register(s) at 0x%X result: 0x%X", quantity, address, res)

	return hex.EncodeToString(res), nil
}
