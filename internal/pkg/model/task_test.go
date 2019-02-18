package model_test

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/input"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/output"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	task := model.NewTask()

	if task.Inputs == nil || len(task.Inputs) != 0 {
		t.Fail()
	}

	if task.Outputs == nil || len(task.Outputs) != 0 {
		t.Fail()
	}
}

type TestInput struct {
	inputChannel chan<- model.Input
}

func (t *TestInput) SetConfig(c model.Config) bool {
	return true
}

func (t *TestInput) GetConfig() model.Config {
	return nil
}

func (t *TestInput) SetInputChannel(c chan<- model.Input) {
	t.inputChannel = c
}

func (t *TestInput) Listen() {
	time.Sleep(2 * time.Second)
	t.inputChannel <- t
}

func (t *TestInput) Stop() {

}

type TestOutput struct {
	C chan<- bool
}

func (*TestOutput) SetConfig(c model.Config) bool {
	panic("implement me")
}

func (t *TestOutput) GetConfig() model.Config {
	return nil
}

func (t *TestOutput) Execute() bool {
	if t.C != nil {
		t.C <- true
	}

	return true
}

func TestTask_RegisterInput(t *testing.T) {
	task := model.NewTask()
	in := &TestInput{}

	task.RegisterInput(in)
	i := task.Inputs[0]

	if i != in {
		t.Fail()
	}
}

func TestTask_RegisterOutput(t *testing.T) {
	task := model.NewTask()
	ou := &TestOutput{}

	task.RegisterOutput(ou)
	o := task.Outputs[0]

	if o != ou {
		t.Fail()
	}
}

func TestTask_Run(t *testing.T) {
	in := &TestInput{}
	ou := &TestOutput{}

	c := make(chan bool)
	ou.C = c

	task := model.NewTask()
	task.RegisterInput(in)
	task.RegisterOutput(ou)

	go task.Run()
	select {
	case res := <-c:
		if !res {
			t.Fail()
		}
	case <-time.After(3 * time.Second):
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_Correct(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task1.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError != nil {
		logrus.Error(yamlError)
	}

	inputConfig := task.Inputs[0].GetConfig()
	if inputConfig.GetName() != "cron input" {
		t.Fail()
	}

	if inputConfig.GetType() != input.CronInputType {
		t.Fail()
	}

	if inputConfig.GetUserConfig()["schedule"] != "* * * * * *" {
		t.Fail()
	}

	outputConfig := task.Outputs[0].GetConfig()
	if outputConfig.GetName() != "HTTP GET google.com" {
		t.Fail()
	}

	if outputConfig.GetType() != output.HTTPOutputType {
		t.Fail()
	}

	if outputConfig.GetUserConfig()["method"] != "GET" {
		t.Fail()
	}

	if outputConfig.GetUserConfig()["url"] != "https://google.com" {
		t.Fail()
	}

	if outputConfig.GetUserConfig()["body"] != "" {
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_YAML_Error(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task_not_yaml.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError == nil {
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_No_Inputs(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task_no_inputs.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError == nil {
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_No_Outputs(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task_no_outputs.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError == nil {
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_Malformed_Inputs(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task_malformed_inputs.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError == nil {
		t.Fail()
	}
}

func TestTask_UnmarshalYAML_Malformed_Outputs(t *testing.T) {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

	data, err := ioutil.ReadFile("test_data/test_task_malformed_outputs.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError == nil {
		t.Fail()
	}
}
