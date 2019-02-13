package model

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	task := NewTask()

	if task.Inputs == nil || len(task.Inputs) != 0 {
		t.Fail()
	}

	if task.Outputs == nil || len(task.Outputs) != 0 {
		t.Fail()
	}
}

type TestInput struct {
	inputChannel chan <- Input
}

func (t *TestInput) UseConfig(c InputConfig) bool {
	return true
}

func (t *TestInput) SetInputChannel(c chan <-Input) {
	t.inputChannel = c
}

func (t *TestInput) Listen() {
	time.Sleep(2 * time.Second)
	t.inputChannel <- t
}

func (t *TestInput) Stop() {

}

type TestOutput struct {
	C chan <- bool
}

func (*TestOutput) UseConfig(c OutputConfig) bool {
	panic("implement me")
}

func (t *TestOutput) Execute() bool {
	if t.C != nil {
		t.C <- true
	}

	return true
}

func TestTask_RegisterInput(t *testing.T) {
	task := NewTask()
	input := &TestInput{}

	task.RegisterInput(input)
	i := task.Inputs[0]

	if i != input {
		t.Fail()
	}
}

func TestTask_RegisterOutput(t *testing.T) {
	task := NewTask()
	output := &TestOutput{}

	task.RegisterOutput(output)
	o := task.Outputs[0]

	if o != output {
		t.Fail()
	}
}

func TestTask_Run(t *testing.T) {
	input := &TestInput{}
	output := &TestOutput{}

	c := make(chan bool)
	output.C = c

	task := NewTask()
	task.RegisterInput(input)
	task.RegisterOutput(output)

	go task.Run()
	select {
	case res := <- c:
		if !res {
			t.Fail()
		}
	case <- time.After(3 * time.Second):
		t.Fail()
	}
}

func TestTask_Unmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("internal/pkg/model/test_data/test_task1.yml")
	if err != nil {
		t.Fail()
	}

	task := NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError != nil {
		t.Fail()
	}

	logrus.Info(task)
}
