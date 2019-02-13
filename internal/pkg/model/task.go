package model

import (
	"github.com/sirupsen/logrus"
)
var logger *logrus.Entry


type Task struct {
	Inputs []Input
	Outputs []Output
	inputChannel chan Input
}

func NewTask() *Task {
	channel := make(chan Input)
	task := &Task {
		make([]Input, 0),
		make([]Output, 0),
		channel,
	}
	logger = logrus.WithField("task", task)

	return task
}

func (t *Task) Run() {
	for _, input := range t.Inputs {
		go input.Listen()
	}

	for info := range t.inputChannel {
		logger.WithFields(logrus.Fields{"input": info}).Debug("Received input")

		for _, output := range t.Outputs {
			go output.Execute()
		}
	}
}

func (t *Task) RegisterInput(input Input) {
	t.Inputs = append(t.Inputs, input)
	input.SetInputChannel(t.inputChannel)
}

func (t *Task) RegisterOutput(output Output) {
	t.Outputs = append(t.Outputs, output)
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})

	if err := unmarshal(m); err != nil {
		logger.WithField("error", err).Error("Cannot unmarshal Task object")
		return err
	}

	logger.Info(m)

	// TODO: Implement YAML logic here

	return nil
}
