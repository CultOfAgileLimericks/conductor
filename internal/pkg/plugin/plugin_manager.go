package plugin

import (
	"reflect"
)

var Manager = newManager()

type Mapping struct {
	Type reflect.Type
	Config reflect.Type
}

type M struct {
	Inputs map[string]Mapping
	Outputs map[string]Mapping
}

func newManager() *M {
	m := &M{}
	m.Inputs = make(map[string]Mapping)
	m.Outputs = make(map[string]Mapping)

	return m
}

func (m *M) RegisterInput(name string, input reflect.Type, config reflect.Type) {
	m.Inputs[name] = Mapping{
		input,
		config,
	}
}

func (m *M) RegisterOutput(name string, output reflect.Type, config reflect.Type) {
	m.Outputs[name] = Mapping{
		output,
		config,
	}
}
