package model

type InputConfig interface {
	InputType() string

	InputName() string
	SetInputName(name string)

	InputUserConfig() map[string]interface{}
	SetInputUserConfig(c map[string]interface{})
}
