package main

import (
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/input"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/output"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

func InitPluginManager() {
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))

}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	InitPluginManager()

	data, err := ioutil.ReadFile("internal/pkg/model/test_data/test_task1.yml")
	if err != nil {
		logrus.Error(err)
	}

	task := model.NewTask()

	yamlError := yaml.Unmarshal(data, task)
	if yamlError != nil {
		logrus.Error(yamlError)
	}

	logrus.Info(task)

	task.Run()

	fmt.Printf("%+v", task)
}
