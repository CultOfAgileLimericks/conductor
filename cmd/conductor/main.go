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
	plugin.Manager.RegisterInput("cron", input.NewCronInput, reflect.TypeOf(input.CronInputConfig{}))
	plugin.Manager.RegisterInput("http", input.NewHTTPInput, reflect.TypeOf(input.HTTPInputConfig{}))
	plugin.Manager.RegisterOutput("http", output.NewHTTPOutput, reflect.TypeOf(output.HTTPOutputConfig{}))
}

func main() {
	InitPluginManager()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
