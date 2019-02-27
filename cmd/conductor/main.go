package main

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/input"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/output"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
)

func InitPluginManager() {
	plugin.Manager.RegisterInput(input.CronInputType, input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput(input.HTTPInputType, input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterInput(input.AMQPInputType, input.NewAMQPInput, reflect.TypeOf(input.AMQPInputConfig{}))
	plugin.Manager.RegisterOutput(output.HTTPOutputType, output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))
}

func main() {
	InitPluginManager()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
