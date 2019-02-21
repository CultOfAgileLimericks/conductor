package model

type Output interface {
	SetConfig(c Config) bool
	GetConfig() Config

	Execute() bool
}
