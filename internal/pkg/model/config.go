package model

type Config interface {
	GetType() string

	GetName() string
	SetName(name string)

	GetUserConfig() map[string]interface{}
	SetUserConfig(c map[string]interface{})
}
