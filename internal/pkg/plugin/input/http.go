package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

var httpInputLogger *logrus.Entry
const HTTP_INPUT_TYPE = "http"

type HTTPInput struct {
	Config model.InputConfig
	channel chan <-model.Input
	server *http.Server
}

type HTTPInputConfig struct {
	Name string
	Addr string
}

func (i *HTTPInputConfig) InputType() string {
	return "http"
}

func (i *HTTPInputConfig) InputName() string {
	return i.Name
}

func (i *HTTPInputConfig) SetInputName(n string) {
	i.Name = n
}

func (i *HTTPInputConfig) InputUserConfig() map[string]interface{} {
	userConfig := make(map[string]interface{})

	userConfig["addr"] = i.Addr
	userConfig["name"] = i.Name

	return userConfig
}

func (i *HTTPInputConfig) SetInputUserConfig(c map[string]interface{})  {
	addr, ok := c["addr"].(string)
	if !ok {
		httpInputLogger.Fatal("addr field not found or incorrect type")
	}
	i.Addr = addr

	name, ok := c["name"].(string)
	if !ok {
		httpInputLogger.Fatal("name field not found or incorrect type")
	}
	i.Name = name
}

func NewHTTPInput() *HTTPInput {
	httpInput :=  &HTTPInput{
		nil,
		nil,
		nil,
	}

	httpInputLogger = logrus.WithField("input", httpInput)

	return httpInput
}

func (input *HTTPInput) UseConfig(c model.InputConfig) bool {
	if c.InputName() == "" || c.InputType() != HTTP_INPUT_TYPE {
		return false
	}

	if c := c.InputUserConfig(); c == nil {
		return false
	}

	input.Config = c
	return true
}

func (input *HTTPInput) SetInputChannel(c chan <-model.Input) {
	input.channel = c
}

func (input *HTTPInput) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- input
	})

	httpInputConfig := input.Config.(*HTTPInputConfig)
	input.server = &http.Server{
		Addr: httpInputConfig.Addr,
	}

	if err := input.server.ListenAndServe(); err != http.ErrServerClosed {
		// TODO: Possible race condition, implement better error handling
		httpInputLogger.WithField("error", err).Error("Failed to start server")
	}
}

func (input *HTTPInput) Stop() {
	if err := input.server.Shutdown(nil); err != nil {
		httpInputLogger.WithField("error", err).Error("Cannot shut down server")
	}
}
