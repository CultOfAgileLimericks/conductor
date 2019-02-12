package model

type Output interface {
	UseConfig(c OutputConfig) bool


	Execute() bool
}