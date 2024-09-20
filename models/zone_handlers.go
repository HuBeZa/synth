package models

import (
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type ZoneHandler[Model tea.Model] func(Model, tea.MouseMsg) (tea.Model, tea.Cmd)
type ZoneHandlers[Model tea.Model] map[string]ZoneHandler[Model]

func (zoneHandlers ZoneHandlers[Model]) Handle(m Model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	for id, handler := range zoneHandlers {
		if zone.Get(id).InBounds(msg) {
			return handler(m, msg)
		}
	}
	return m, nil
}
