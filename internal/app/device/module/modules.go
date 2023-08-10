package module

import (
	"sync"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/errors"
)

type Modules struct {
	id      uint64
	modules []Module
	mu      sync.Mutex
}

type Entity struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Type Type   `json:"type"`
}

func (m *Modules) ReadAll() []*Entity {
	m.mu.Lock()
	defer m.mu.Unlock()

	list := make([]*Entity, 0, len(m.modules))
	for _, mod := range m.modules {
		list = append(list, mod.Entity())
	}

	return list
}

func (m *Modules) CreateOne(kind Type, config Config) (Module, error) {
	mod, err := newModule(&m.id, kind, config)
	if err != nil {
		return nil, errors.Wrap(err, "create module")
	}

	return m.add(mod), nil
}

func (m *Modules) add(mod Module) Module {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.modules = append(m.modules, mod)
	return mod
}

func (m *Modules) ReadOne(id uint64) (Module, error) {
	idx, mod := m.find(id)
	if idx == -1 {
		return nil, errors.Wrap(&ErrModuleNotFound{ID: id}, "read module")
	}

	return mod, nil
}

func (m *Modules) find(id uint64) (int, Module) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for idx, mod := range m.modules {
		if id == mod.ID() {
			return idx, mod
		}
	}

	return -1, nil
}

type Config struct {
	Name string `json:"name"`
}

func (m *Modules) UpdateOne(id uint64, config Config) error {
	mod, err := m.ReadOne(id)
	if err != nil {
		return errors.Wrap(err, "update module")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	mod.Update(config)
	return nil
}

func (m *Modules) DeleteOne(id uint64) error {
	idx, _ := m.find(id)
	if idx == -1 {
		return errors.Wrap(&ErrModuleNotFound{ID: id}, "delete module")
	}

	m.delete(idx)
	return nil
}

func (m *Modules) delete(index int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.modules = append(m.modules[:index], m.modules[index+1:]...)
}
