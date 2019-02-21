package model

type Input interface {
	SetConfig(c Config) bool
	GetConfig() Config

	SetInputChannel(c chan<- Input)
	Listen()
	Stop()
}
