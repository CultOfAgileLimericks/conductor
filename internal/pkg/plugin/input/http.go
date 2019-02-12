package input

import (
	"context"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

var httpInputLogger *logrus.Entry
const HTTPInputType = "http"

type HTTPInput struct {
	Config model.InputConfig
	channel chan <-model.Input
	err chan error
	server *http.Server
}

type HTTPInputConfig struct {
	Name string
	Addr string
}

func (i *HTTPInputConfig) InputType() string {
	return HTTPInputType
}

func (i *HTTPInputConfig) InputName() string {
	return i.Name
}

func (i *HTTPInputConfig) SetInputName(n string) {
	i.Name = n
}

func (i *HTTPInputConfig) InputUserConfig() map[string]interface{} {
	if i.Addr == "" {
		return nil
	}

	userConfig := make(map[string]interface{})

	userConfig["addr"] = i.Addr

	return userConfig
}

func (i *HTTPInputConfig) SetInputUserConfig(c map[string]interface{})  {
	logEntry := logrus.WithField("config", i)
	addr, ok := c["addr"].(string)
	if !ok {
		logEntry.Error("addr field not found or incorrect type")
	}
	i.Addr = addr
}

func NewHTTPInput() *HTTPInput {
	httpInput :=  &HTTPInput{
		nil,
		nil,
		make(chan error),
		nil,
	}

	httpInputLogger = logrus.WithField("input", httpInput)

	return httpInput
}

func (input *HTTPInput) UseConfig(c model.InputConfig) bool {
	if c.InputName() == "" || c.InputType() != HTTPInputType {
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
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- input
	})

	httpInputConfig := input.Config.(*HTTPInputConfig)
	input.server = &http.Server{
		Addr: httpInputConfig.Addr,
		Handler: handler,
	}

	go func() {
		err := input.server.ListenAndServe()
		input.err <- err

	}()

	select {
	case err := <-input.err:
		if err != http.ErrServerClosed {
			httpInputLogger.WithField("error", err).Error("Failed to start server")
		}
		close(input.err)
	}
}

func (input *HTTPInput) Stop() {
	if err := input.server.Shutdown(context.Background()); err != nil {
		httpInputLogger.WithField("error", err).Error("Cannot shut down server")
	}
}
