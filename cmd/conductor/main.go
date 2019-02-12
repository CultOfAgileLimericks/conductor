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
	//	Name:"listen on :8080",
	//	Addr: ":8080",
	//
	//}
	//
	//httpInput.UseConfig(inputConfig)

	cronInput := input.NewCronInput()
	inputConfig := &input.CronInputConfig{
		Name: "run every minute",
		Schedule: "0 * * * * *",
	}
	cronInput.UseConfig(inputConfig)

	httpOutput := output.NewHTTPOutput()
	outputConfig := &output.HTTPOutputConfig{
		Name: "GET bing.com",
		Method: "GET",
		URL: "https://bing.com",
		Body: "",
	}
	httpOutput.UseConfig(outputConfig)

	t.RegisterInput(cronInput)
	t.RegisterOutput(httpOutput)

	t.Run()

	fmt.Printf("%+v", t)
}
