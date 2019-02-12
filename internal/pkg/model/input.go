package model

type Input interface {
	UseConfig(c InputConfig)

	SetInputChannel(c chan <-Input)
	Listen()
}
