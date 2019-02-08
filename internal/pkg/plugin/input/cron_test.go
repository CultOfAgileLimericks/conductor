package input

import "testing"

func TestNewCronInput(t *testing.T) {
	c := NewCronInput()
	if c.Config != nil {
		t.Fail()
	}
}
