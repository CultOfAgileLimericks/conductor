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

func InitManager() {
	plugin.Manager.RegisterInput("cron", reflect.TypeOf(input.CronInput{}), reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", reflect.TypeOf(input.HTTPInput{}), reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", reflect.TypeOf(output.HTTPOutput{}), reflect.TypeOf(output.HTTPOutputConfig{}))

}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	InitManager()

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

	//t := model.NewTask()
	//httpInput := input.NewHTTPInput()
	//inputConfig := &input.HTTPInputConfig{
	//	GetName:"listen on :8080",
	//	Addr: ":8080",
	//
	//}
	//
	//httpInput.SetConfig(inputConfig)

	//cronInput := input.NewCronInput()
	//inputConfig := &input.CronInputConfig{
	//	Schedule: "0 * * * * *",
	//}
	//inputConfig.SetName("run every minute")
	//cronInput.SetConfig(inputConfig)
	//
	//httpOutput := output.NewHTTPOutput()
	//outputConfig := &output.HTTPOutputConfig{
	//	Method: "GET",
	//	URL: "https://bing.com",
	//	Body: "",
	//}
	//outputConfig.SetName("GET bing.com")
	//httpOutput.SetConfig(outputConfig)
	//
	//t.RegisterInput(cronInput)
	//t.RegisterOutput(httpOutput)

	task.Run()

	fmt.Printf("%+v", task)
}
