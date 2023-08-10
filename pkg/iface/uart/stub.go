//go:build !sysboard

package uart

type Conn struct{}

func (u *UART) Connect(id byte) *Conn {
	return &Conn{}
}

func (c *Conn) ReadInputRegisters(address, quantity uint16) (string, error) {
	return "", nil
}

func (c *Conn) ReadHoldingRegisters(address, quantity uint16) (results string, err error) {
	return "", nil
}
