package model

type Output interface {
	UseConfig(c OutputConfig)


	Execute() bool
}