package model

type OutputConfig interface {
	OutputType() string

	OutputName() string
	SetOutputName(name string)

	OutputUserConfig() map[string]interface{}
	SetOutputUserConfig(config map[string]interface{})
}
