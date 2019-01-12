package model

type Input interface {
	SetInputChannel(c chan <-[]byte)
	Listen()
}
