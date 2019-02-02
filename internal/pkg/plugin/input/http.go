package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger *logrus.Entry

type HTTPInput struct {
	Config model.InputConfig
	channel chan <-model.Input
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
		logger.Fatal("addr field not found or incorrect type")
	}
	i.Addr = addr

	name, ok := c["name"].(string)
	if !ok {
		logger.Fatal("name field not found or incorrect type")
	}
	i.Name = name
}

func NewHTTPInput() *HTTPInput {
	httpInput :=  &HTTPInput{
		nil,
		nil,
	}

	logger = logrus.WithField("input", httpInput)

	return httpInput
}

func (input *HTTPInput) UseConfig(c model.InputConfig) {
	input.Config = c
}

func (input *HTTPInput) SetInputChannel(c chan <-model.Input) {
	input.channel = c
}

func (input *HTTPInput) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- input
	})

	httpInputConfig := input.Config.(*HTTPInputConfig)
	logger.Fatal(http.ListenAndServe(httpInputConfig.Addr, nil))
}
