package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var cronInputLogger *logrus.Entry
const CRON_INPUT_TYPE = "cron"

type CronInput struct {
	Config model.InputConfig
	channel chan <-model.Input
	cron *cron.Cron
	stop chan bool
}

type CronInputConfig struct {
	Name string
	Schedule string
}

func (i *CronInputConfig) InputType() string {
	return CRON_INPUT_TYPE
}

func (i *CronInputConfig) InputName() string {
	return i.Name
}

func (i *CronInputConfig) SetInputName(n string) {
	i.Name = n
}

func (i *CronInputConfig) InputUserConfig() map[string]interface{} {
	userConfig := make(map[string]interface{})

	userConfig["schedule"] = i.Schedule

	return userConfig
}

func (i *CronInputConfig) SetInputUserConfig(c map[string]interface{})  {
	schedule, ok := c["schedule"].(string)
	if !ok {
		cronInputLogger.Fatal("schedule field not found or incorrect type")
	}
	i.Schedule = schedule
}

func NewCronInput() *CronInput {
	cronInput :=  &CronInput{
		nil,
		nil,
		nil,
		make(chan bool),
	}

	cronInputLogger = logrus.WithField("input", cronInput)

	return cronInput
}

func (input *CronInput) UseConfig(c model.InputConfig) bool{
	if c.InputName() == "" || c.InputType() != CRON_INPUT_TYPE {
		return false
	}

	if c := c.InputUserConfig(); c == nil {
		return false
	}

	input.Config = c
	return true
}

func (input *CronInput) SetInputChannel(c chan <-model.Input) {
	input.channel = c
}

func (input *CronInput) Listen() {
	input.cron = cron.New()
	err := input.cron.AddFunc(input.Config.(*CronInputConfig).Schedule, func() {
		input.channel <- input
	})

	if err != nil {
		cronInputLogger.WithField("error", err).Error("Failed to add cron schedule")
	}

	input.cron.Start()
	select {
		case <- input.stop:
			input.cron.Stop()
			return
	}
}

func (input *CronInput) Stop() {
	input.stop <- true
}
