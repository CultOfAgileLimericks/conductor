package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var cronInputLogger *logrus.Entry
const CronInputType = "cron"

type CronInput struct {
	config model.Config
	channel chan <-model.Input
	cron *cron.Cron
	stop chan bool
}

type CronInputConfig struct {
	name     string
	Schedule string
}

func (i *CronInputConfig) GetType() string {
	return CronInputType
}

func (i *CronInputConfig) GetName() string {
	return i.name
}

func (i *CronInputConfig) SetName(n string) {
	i.name = n
}

func (i *CronInputConfig) GetUserConfig() map[string]interface{} {
	if i.Schedule == "" {
		return nil
	}

	userConfig := make(map[string]interface{})

	userConfig["Schedule"] = i.Schedule

	return userConfig
}

func (i *CronInputConfig) SetUserConfig(c map[string]interface{})  {
	logEntry := logrus.WithField("config", i)
	schedule, ok := c["Schedule"].(string)
	if !ok {
		logEntry.Error("Schedule field not found or incorrect type")
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

func (input *CronInput) SetConfig(c model.Config) bool{
	if c.GetName() == "" || c.GetType() != CronInputType {
		return false
	}

	if userConfig := c.GetUserConfig(); userConfig == nil {
		return false
	}

	input.config = c
	return true
}

func (input *CronInput) GetConfig() model.Config {
	return input.config
}

func (input *CronInput) SetInputChannel(c chan <-model.Input) {
	input.channel = c
}

func (input *CronInput) Listen() {
	input.cron = cron.New()
	err := input.cron.AddFunc(input.config.(*CronInputConfig).Schedule, func() {
		input.channel <- input
	})

	if err != nil {
		cronInputLogger.WithField("error", err).Error("Failed to add cron Schedule")
		// TODO: Add error handling like sending error messages over channel
	}

	input.cron.Start()
	select {
		case <- input.stop:
			input.cron.Stop()
	}

	close(input.channel)
}

func (input *CronInput) Stop() {
	input.stop <- true
}
