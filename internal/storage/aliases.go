package storage

import "github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/types"

func (s *Storage) ListAliases() ([]*types.Alias, error) {
	return s.aliasesDB.All().([]*types.Alias), nil
}

func (s *Storage) AddAlias(alias *types.Alias) error {
	id, err := s.aliasesDB.Append(alias)
	if err != nil {
		return err
	}
	alias.ID = id
	return s.UpdateAlias(alias)
}

func (s *Storage) UpdateAlias(alias *types.Alias) error {
	return s.aliasesDB.Set(alias.ID, alias)
}

func (s *Storage) DeleteAlias(id int) error {
	return s.aliasesDB.Delete(id)
}
