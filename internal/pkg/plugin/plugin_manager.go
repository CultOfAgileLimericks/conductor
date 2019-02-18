package plugin

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
)

type Manager struct {
	Inputs map[string]model.Input
	Outputs map[string]model.Output
}

func NewManager() *Manager {
	m := &Manager{}
	m.Inputs = make(map[string]model.Input)
	m.Outputs = make(map[string]model.Output)

	return m
}

func (m *Manager) RegisterInput(i model.Input) {
	m.Inputs[i.GetConfig().GetName()] = i
}

func (m *Manager) RegisterOutput(o model.Output) {
	m.Outputs[o.GetConfig().GetName()] = o
}
