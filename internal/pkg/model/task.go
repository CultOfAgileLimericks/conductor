package model

import (
	"errors"
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin"
	"github.com/sirupsen/logrus"
	"reflect"
)

var logger *logrus.Entry

type Task struct {
	// TODO: Add Name field
	Inputs       []Input
	Outputs      []Output
	inputChannel chan Input
}

func NewTask() *Task {
	channel := make(chan Input)
	task := &Task{
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

	inputs, ok := m["inputs"].([]interface{})
	logger.Info(reflect.TypeOf(m["inputs"]))
	if !ok {
		msg := "'inputs' field not found or unexpected format"
		logger.WithField("inputs", m["inputs"]).Error(msg)
		return errors.New(msg)
	}

	outputs, ok := m["outputs"].([]interface{})
	if !ok {
		msg := "'outputs' field not found or unexpected format"
		logger.WithField("outputs", m["outputs"]).Error(msg)
		return errors.New(msg)
	}

	for _, yamlInput := range inputs {
		mapInput, ok := yamlInput.(map[interface{}]interface{})
		if !ok {
			msg := "'inputs' field not found or unexpected format"
			logger.WithField("inputs", m["inputs"]).Error(msg)
			return errors.New(msg)
		}
		_type := mapInput["type"].(string)
		input := plugin.Manager.Inputs[_type].Type().(Input)
		config := reflect.New(plugin.Manager.Inputs[_type].Config).Interface().(Config)

		c := make(map[string]interface{})

		for key, value := range mapInput["config"].(map[interface{}]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			c[strKey] = value
		}

		config.SetUserConfig(c)
		config.SetName(fmt.Sprintf("%v", c["name"]))
		input.SetConfig(config)

		t.RegisterInput(input)
		logger.Info(reflect.TypeOf(yamlInput))
	}

	for _, yamlOutput := range outputs {
		mapOutput, ok := yamlOutput.(map[interface{}]interface{})
		if !ok {
			msg := "'outputs' field not found or unexpected format"
			logger.WithField("outputs", m["outputs"]).Error(msg)
			return errors.New(msg)
		}
		_type := mapOutput["type"].(string)
		output := plugin.Manager.Outputs[_type].Type().(Output)
		config := reflect.New(plugin.Manager.Outputs[_type].Config).Interface().(Config)

		c := make(map[string]interface{})

		for key, value := range mapOutput["config"].(map[interface{}]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			c[strKey] = value
		}

		config.SetUserConfig(c)
		config.SetName(fmt.Sprintf("%v", c["name"]))
		output.SetConfig(config)

		t.RegisterOutput(output)
	}

	return nil
}
