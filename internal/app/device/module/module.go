package module

import (
	"fmt"
)

type Type string

const (
	TypeDTH          Type = "DTH-XX"
	TypeDS18B20      Type = "DS18B20"
	TypeInfrared     Type = "Инфракрасный"
	TypeResistive    Type = "Резистивный"
	TypeThermocouple Type = "Термопара на MLX90614"
	TypeTriac2       Type = "Triac-2"
	TypeTriac4       Type = "Triac-4"
	TypeTriac6       Type = "Triac-6"
	TypeRelay2       Type = "Relay-2"
	TypeRelay4       Type = "Relay-4"
	TypeRelay6       Type = "Relay-6"
)

type ErrUnknownType struct {
	Kind Type
}

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("unknown module type %q", e.Kind)
}

type ErrModuleNotFound struct {
	ID uint64
}

func (e *ErrModuleNotFound) Error() string {
	return fmt.Sprintf("module with id %d not found", e.ID)
}

type Module interface {
	ID() uint64
	Entity() *Entity
	Update(config Config)
}

func newModule(id *uint64, kind Type, config Config) (Module, error) {
	builder, ok := map[Type]func(id *uint64, config Config) Module{
		TypeDTH:          newDHT,
		TypeDS18B20:      newDS18B20,
		TypeInfrared:     newInfrared,
		TypeResistive:    newResistive,
		TypeThermocouple: newThermocouple,
		TypeTriac2:       newTriac2,
		TypeTriac4:       newTriac4,
		TypeTriac6:       newTriac6,
		TypeRelay2:       newRelay2,
		TypeRelay4:       newRelay4,
		TypeRelay6:       newRelay6,
	}[kind]
	if !ok {
		return nil, &ErrUnknownType{Kind: kind}
	}

	return builder(id, config), nil
}

type module struct {
	id   uint64
	name string
	kind Type
}

func (m *module) ID() uint64 {
	return m.id
}

func (m *module) Entity() *Entity {
	return &Entity{
		ID:   m.id,
		Name: m.name,
		Type: m.kind,
	}
}

func (m *module) Update(config Config) {
	m.name = config.Name
}
