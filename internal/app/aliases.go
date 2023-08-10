package app

import (
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/alias"
)

func (a *Application) AliasesList(req struct {
	Page  int `json:"page" validate:"gte=1"`
	Limit int `json:"limit" validate:"oneof=1 2 3"`
}) ([]*alias.Alias, error) {
	return []*alias.Alias{
		{
			ID:    0,
			Name:  "Температура Офис",
			State: "22 ℃",
			Path:  "MK -> Temp Resistive",
		},
		{
			ID:    1, //nolint:gomnd
			Name:  "Температура Теплица",
			State: "18 ℃",
			Path:  "MK -> Temp Resistive",
		},
		{
			ID:    2, //nolint:gomnd
			Name:  "Реле отопления",
			State: "ON",
			Path:  "MK -> Relay-4 -> Out-1",
		},
	}, nil
}

func (a *Application) AliasesCount() (int, error) {
	return 0, nil
}
