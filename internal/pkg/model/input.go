package model

type Input interface {
	SetInputChannel(c chan <-Input)
	Listen()
}
