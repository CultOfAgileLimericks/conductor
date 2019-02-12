package input

import "testing"

func TestNewCronInput(t *testing.T) {
	c := NewCronInput()
	if c.Config != nil {
		t.Fail()
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