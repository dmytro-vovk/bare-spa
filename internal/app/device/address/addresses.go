package address

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type Address uint

type Borders struct {
	Left  Address
	Right Address
}

var addresses *Addresses

type Addresses struct {
	borders   Borders
	addresses []Address
	mu        sync.Mutex
}

const (
	LeftBorder  = 1
	RightBorder = 1<<8 - 2
)

func newAddresses(borders Borders) (*Addresses, error) {
	if borders.Left > borders.Right {
		return nil, &ErrWrongBorders{
			LeftBorder:  borders.Left,
			RightBorder: borders.Right,
		}
	}

	slice := make([]Address, 0, borders.Right-borders.Left)
	for addr := borders.Left; addr <= borders.Right; addr++ {
		slice = append(slice, addr)
	}

	return &Addresses{
		borders:   borders,
		addresses: slice,
	}, nil
}

func init() {
	var err error
	addresses, err = newAddresses(Borders{
		Left:  LeftBorder,
		Right: RightBorder,
	})
	if err != nil {
		panic(errors.Wrap(err, "could not init addresses"))
	}
}

type ErrWrongBorders struct {
	LeftBorder  Address
	RightBorder Address
}

func (e *ErrWrongBorders) Error() string {
	return fmt.Sprintf("first address 0x%X greater than last address 0x%X", e.LeftBorder, e.RightBorder)
}

func ReadAll() []Address {
	return addresses.ReadAll()
}

func (a *Addresses) ReadAll() []Address {
	a.mu.Lock()
	defer a.mu.Unlock()

	clone := make([]Address, len(a.addresses))
	copy(clone, a.addresses)

	return clone
}

type ErrAddressOutOfRange struct {
	Address Address
}

func (e *ErrAddressOutOfRange) Error() string {
	return fmt.Sprintf("address 0x%X out of range", e.Address)
}

func (a *Addresses) check(address Address) error {
	if address < a.borders.Left || address > a.borders.Right {
		return errors.Wrap(&ErrAddressOutOfRange{Address: address}, "check address")
	}

	return nil
}

func (a *Addresses) findIndex(address Address) int {
	a.mu.Lock()
	defer a.mu.Unlock()

	for idx, addr := range a.addresses {
		if addr == address {
			return idx
		}
	}

	return -1
}

type ErrAddressExist struct {
	Address Address
}

func (e *ErrAddressExist) Error() string {
	return fmt.Sprintf("address 0x%X already exist", e.Address)
}

func Release(address Address) error {
	return addresses.Release(address)
}

func (a *Addresses) Release(address Address) error {
	const errFormat = "release address"

	if err := a.check(address); err != nil {
		return errors.Wrap(err, errFormat)
	}

	idx := a.findIndex(address)
	if idx != -1 {
		return errors.Wrap(&ErrAddressExist{Address: address}, errFormat)
	}

	a.insert(address)
	return nil
}

func (a *Addresses) insert(address Address) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for idx, addr := range a.addresses {
		if addr > address {
			a.addresses = append(a.addresses[:idx+1], a.addresses[idx:]...)
			a.addresses[idx] = address
			return
		}
	}

	a.addresses = append(a.addresses, address)
}

type ErrAddressNotFound struct {
	Address Address
}

func (e *ErrAddressNotFound) Error() string {
	return fmt.Sprintf("address 0x%X not found", e.Address)
}

func Borrow(address Address) error {
	return addresses.Borrow(address)
}

func (a *Addresses) Borrow(address Address) error {
	const errFormat = "borrow address"

	if err := a.check(address); err != nil {
		return errors.Wrap(err, errFormat)
	}

	idx := a.findIndex(address)
	if idx == -1 {
		return errors.Wrap(&ErrAddressNotFound{Address: address}, errFormat)
	}

	a.delete(idx)
	return nil
}

func (a *Addresses) delete(index int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.addresses = append(a.addresses[:index], a.addresses[index+1:]...)
}
