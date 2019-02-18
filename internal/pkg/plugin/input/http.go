package input

import (
	"context"
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var httpInputLogger *logrus.Entry

const HTTPInputType = "http"

type HTTPInput struct {
	config  model.Config
	channel chan<- model.Input
	err     chan error
	server  *http.Server
	mutex   *sync.Mutex
}

type HTTPInputConfig struct {
	name string
	Addr string
}

func (i *HTTPInputConfig) GetType() string {
	return HTTPInputType
}

func (i *HTTPInputConfig) GetName() string {
	return i.name
}

func (i *HTTPInputConfig) SetName(n string) {
	i.name = n
}

func (i *HTTPInputConfig) GetUserConfig() map[string]interface{} {
	if i.Addr == "" {
		return nil
	}

	userConfig := make(map[string]interface{})

	userConfig["addr"] = i.Addr

	return userConfig
}

func (i *HTTPInputConfig) SetUserConfig(c map[string]interface{}) {
	logEntry := logrus.WithField("config", i)
	if c["addr"] == nil {
		logEntry.Error("addr field not found or incorrect type")
	} else {
		addr := fmt.Sprintf("%v", c["addr"])
		i.Addr = addr
	}
}

func NewHTTPInput() interface{} {
	httpInput := &HTTPInput{
		nil,
		nil,
		make(chan error),
		nil,
		&sync.Mutex{},
	}

	httpInputLogger = logrus.WithField("input", httpInput)

	return httpInput
}

func (input *HTTPInput) SetConfig(c model.Config) bool {
	if c.GetName() == "" || c.GetType() != HTTPInputType {
		return false
	}

	if c := c.GetUserConfig(); c == nil {
		return false
	}

	input.config = c
	return true
}

func (input *HTTPInput) GetConfig() model.Config {
	return input.config
}

func (input *HTTPInput) SetInputChannel(c chan<- model.Input) {
	input.channel = c
}

func (input *HTTPInput) Listen() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- input
	})

	httpInputConfig := input.config.(*HTTPInputConfig)

	input.mutex.Lock()
	input.server = &http.Server{
		Addr:    httpInputConfig.Addr,
		Handler: handler,
	}
	input.mutex.Unlock()

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

		return
	}
}

func (input *HTTPInput) Stop() {
	input.mutex.Lock()
	if err := input.server.Shutdown(context.Background()); err != nil {
		httpInputLogger.WithField("error", err).Error("Cannot shut down server")
	}
	input.mutex.Unlock()
}
