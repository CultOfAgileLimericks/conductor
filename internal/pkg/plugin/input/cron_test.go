package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"testing"
	"time"
)

func TestNewCronInput(t *testing.T) {
	c := NewCronInput().(*CronInput)
	if c.GetConfig() != nil {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputName(t *testing.T) {
	config := &CronInputConfig{}
	config.SetName("test name")

	if config.GetName() != "test name" {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputUserConfig_Correct(t *testing.T) {
	config := &CronInputConfig{}
	userConfig := make(map[string]interface{})
	userConfig["schedule"] = "* * * * * *"
	config.SetUserConfig(userConfig)

	if config.GetUserConfig()["schedule"] != userConfig["schedule"] {
		t.Fail()
	}

	if config.Schedule != userConfig["schedule"] {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputUserConfig_Incorrect(t *testing.T) {
	config := &CronInputConfig{}
	userConfig := make(map[string]interface{})
	config.SetUserConfig(userConfig)

	if config.Schedule != "" {
		t.Fail()
	}

	if config.GetUserConfig()["schedule"] != nil {
		t.Fail()
	}
}

func TestCronInput_UseConfig_Incorrect(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{}

	ok := c.SetConfig(config)
	if ok || c.config == config {
		t.Fail()
	}

	config.SetName("TestCronInput")

	ok = c.SetConfig(config)
	if ok || c.GetConfig() == config {
		t.Fail()
	}
}

func TestCronInput_UseConfig_Correct(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{
		Schedule: "* * * * * *",
	}
	config.SetName("TestCronInput")

	ok := c.SetConfig(config)
	if !ok || c.GetConfig() != config {
		t.Fail()
	}
}

func TestCronInput_SetInputChannel(t *testing.T) {
	c := NewCronInput().(*CronInput)

	channel := make(chan model.Input)
	c.SetInputChannel(channel)

	if c.channel != channel {
		t.Fail()
	}
}

func TestCronInput_Listen_Correct(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{
		Schedule: "* * * * * *",
		name:     "TestCronInput",
	}

	c.SetConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	defer c.Stop()

	select {
	case i := <-channel:
		if i != c {
			t.Fail()
		}
	case <-time.After(2 * time.Second):
		t.Fail()
	}
}

func TestCronInput_Listen_Incorrect(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{
		Schedule: "asdfgh",
	}
	config.SetName("TestCronInput")

	c.SetConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	defer c.Stop()

	select {
	case i := <-channel:
		if i != c {
			t.Fail()
		}
	case <-time.After(2 * time.Second):

	}
}

func TestCronInput_Listen_Stopped(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{
		Schedule: "* * * * * *",
	}
	config.SetName("TestCronInput")

	c.SetConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	seconds := 0

	for {
		select {
		case i, ok := <-channel:
			if seconds == 3 {
				c.Stop()
			}
			if !ok {
				return
			}
			if i != c {
				t.Fail()
				return
			}
			seconds++
		case <-time.After(4 * time.Second):
			t.Fail()
			return
		}
	}

}

func TestCronInput_UseConfig(t *testing.T) {
	c := NewCronInput().(*CronInput)
	config := &CronInputConfig{}

	c.SetConfig(config)
	if c.GetConfig() == config {
		t.Fail()
	}
}
