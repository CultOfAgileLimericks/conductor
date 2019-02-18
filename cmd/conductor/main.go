package main

import (
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/input"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/output"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	t := model.NewTask()
	//httpInput := input.NewHTTPInput()
	//inputConfig := &input.HTTPInputConfig{
	//	GetName:"listen on :8080",
	//	Addr: ":8080",
	//
	//}
	//
	//httpInput.SetConfig(inputConfig)

	cronInput := input.NewCronInput()
	inputConfig := &input.CronInputConfig{
		Schedule: "0 * * * * *",
	}
	inputConfig.SetName("run every minute")
	cronInput.SetConfig(inputConfig)

	httpOutput := output.NewHTTPOutput()
	outputConfig := &output.HTTPOutputConfig{
		Method: "GET",
		URL: "https://bing.com",
		Body: "",
	}
	outputConfig.SetName("GET bing.com")
	httpOutput.SetConfig(outputConfig)

	t.RegisterInput(cronInput)
	t.RegisterOutput(httpOutput)

	t.Run()

	fmt.Printf("%+v", t)
}
