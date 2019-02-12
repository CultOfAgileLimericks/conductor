package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPInput(t *testing.T) {
	h := NewHTTPInput()
	if h.Config != nil {
		t.Fail()
	}
}

func TestHTTPInputConfig_SetInputName(t *testing.T) {
	c := &HTTPInputConfig{}

	c.SetInputName("test")
	if c.Name != "test" {
		t.Fail()
	}
}

func TestHTTPInputConfig_SetInputUserConfig_Correct(t *testing.T) {
	config := &HTTPInputConfig{}
	userConfig := make(map[string]interface{})
	userConfig["addr"] = "localhost:8080"
	config.SetInputUserConfig(userConfig)

	if config.InputUserConfig()["addr"] != userConfig["addr"] {
		t.Fail()
	}

	if config.Addr != userConfig["addr"] {
		t.Fail()
	}
}

func TestHTTPInputConfig_SetInputUserConfig_Incorrect(t *testing.T) {
	config := &HTTPInputConfig{}
	userConfig := make(map[string]interface{})
	config.SetInputUserConfig(userConfig)

	if config.Addr != "" {
		t.Fail()
	}

	if config.InputUserConfig()["addr"] != nil {
		t.Fail()
	}
}

func TestHTTPInput_UseConfig_Correct(t *testing.T) {
	h := NewHTTPInput()
	config := &HTTPInputConfig{
		Name: "test",
		Addr: "localhost:8080",
	}

	ok := h.UseConfig(config)

	if !ok || h.Config != config {
		t.Fail()
	}
}

func TestHTTPInput_UseConfig_Incorrect(t *testing.T) {
	h := NewHTTPInput()
	config := &HTTPInputConfig{}

	ok := h.UseConfig(config)

	if ok || h.Config != nil {
		t.Fail()
	}

	config.Name = "test name"

	ok = h.UseConfig(config)

	if ok || h.Config != nil {
		t.Fail()
	}
}

func TestHTTPInput_SetInputChannel(t *testing.T) {
	c := NewHTTPInput()

	channel := make(chan model.Input)
	c.SetInputChannel(channel)

	if c.channel != channel {
		t.Fail()
	}
}

func TestHTTPInput_Listen_Correct(t *testing.T) {
	h := NewHTTPInput()
	config := &HTTPInputConfig{
		Name: "test http input",
		Addr: "localhost:8080",
	}

	h.UseConfig(config)
	channel := make(chan model.Input)

	h.SetInputChannel(channel)

	go h.Listen()
	defer h.Stop()

	go func() {
		time.Sleep(1 * time.Second)
		_, err := http.Get("http://localhost:8080")
		if err != nil {
			t.Fail()
			logrus.Info(err)
		}
	}()

	select {
	case i := <- channel:
		if i != h {
			t.Fail()
		}
	case <- time.After(2 * time.Second):
		t.Fail()
	}
}

func TestHTTPInput_Listen_Incorrect(t *testing.T) {
	h := NewHTTPInput()
	config := &HTTPInputConfig{
		Name: "test http input",
		Addr: "localhost:22",
	}

	h.UseConfig(config)
	channel := make(chan model.Input)

	go h.Listen()
	defer h.Stop()

	select {
	case i := <- channel:
		if i == h {
			t.Fail()
		}
	case <- time.After(2 * time.Second):
	}
}
