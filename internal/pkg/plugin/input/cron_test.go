package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"testing"
	"time"
)

func TestNewCronInput(t *testing.T) {
	c := NewCronInput()
	if c.Config != nil {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputName(t *testing.T) {
	config := &CronInputConfig{}
	config.SetInputName("test name")

	if config.Name != "test name" {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputUserConfig_Correct(t *testing.T) {
	config := &CronInputConfig{}
	userConfig := make(map[string]interface{})
	userConfig["schedule"] = "* * * * * *"
	config.SetInputUserConfig(userConfig)

	if config.InputUserConfig()["schedule"] != userConfig["schedule"] {
		t.Fail()
	}

	if config.Schedule != userConfig["schedule"] {
		t.Fail()
	}
}

func TestCronInputConfig_SetInputUserConfig_Incorrect(t *testing.T) {
	config := &CronInputConfig{}
	userConfig := make(map[string]interface{})
	config.SetInputUserConfig(userConfig)

	if config.Schedule != "" {
		t.Fail()
	}
}

func TestCronInput_UseConfig_Incorrect(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{}

	ok := c.UseConfig(config)
	if ok || c.Config == config {
		t.Fail()
	}

	config.Name = "TestCronInput"

	ok = c.UseConfig(config)
	if ok || c.Config == config {
		t.Fail()
	}
}

func TestCronInput_UseConfig_Correct(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{
		Schedule: "* * * * * *",
		Name: "TestCronInput",
	}

	ok := c.UseConfig(config)
	if !ok || c.Config != config {
		t.Fail()
	}
}

func TestCronInput_SetInputChannel(t *testing.T) {
	c := NewCronInput()

	channel := make(chan model.Input)
	c.SetInputChannel(channel)

	if c.channel != channel {
		t.Fail()
	}
}

func TestCronInput_Listen_Correct(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{
		Schedule: "* * * * * *",
		Name: "TestCronInput",
	}

	c.UseConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	defer c.Stop()

	select {
	case i := <- channel:
		if i != c {
			t.Fail()
		}
	case <- time.After(2 * time.Second):
		t.Fail()
	}
}

func TestCronInput_Listen_Incorrect(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{
		Schedule: "asdfgh",
		Name: "TestCronInput",
	}

	c.UseConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	defer c.Stop()

	select {
	case i := <- channel:
		if i != c {
			t.Fail()
		}
	case <- time.After(2 * time.Second):

	}
}

func TestCronInput_Listen_Stopped(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{
		Schedule: "* * * * * *",
		Name: "TestCronInput",
	}

	c.UseConfig(config)
	channel := make(chan model.Input)

	c.SetInputChannel(channel)
	go c.Listen()

	seconds := 0

	for {
		select {
		case i, ok := <- channel:
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
		case <- time.After(4 * time.Second):
			t.Fail()
			return
		}
	}

}

func TestCronInput_UseConfig(t *testing.T) {
	c := NewCronInput()
	config := &CronInputConfig{}

	c.UseConfig(config)
	if c.Config == config {
		t.Fail()
	}
}
