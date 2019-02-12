package model

type Input interface {
	UseConfig(c InputConfig) bool

	SetInputChannel(c chan <-Input)
	Listen()
	Stop()
}
