package output

import (
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var httpOutputLogger *logrus.Entry

const HTTPOutputType = "http"

type HTTPOutput struct {
	config     model.Config
	channel    chan<- []byte
	httpClient *http.Client
}

type HTTPOutputConfig struct {
	name   string
	Method string
	URL    string
	Body   string
}

func (c *HTTPOutputConfig) GetType() string {
	return HTTPOutputType
}

func (c *HTTPOutputConfig) GetName() string {
	return c.name
}

func (c *HTTPOutputConfig) SetName(name string) {
	c.name = name
}

func (c *HTTPOutputConfig) GetUserConfig() map[string]interface{} {
	if c.URL == "" || c.Method == "" {
		return nil
	}

	userConfig := make(map[string]interface{})

	userConfig["method"] = c.Method
	userConfig["url"] = c.URL
	userConfig["body"] = c.Body

	return userConfig
}

func (c *HTTPOutputConfig) SetUserConfig(config map[string]interface{}) {
	logEntry := logrus.WithField("config", c)
	if config["method"] == nil {
		logEntry.Error("method field not found or incorrect type")
	} else {
		method := fmt.Sprintf("%v", config["method"])
		c.Method = method
	}

	if config["url"] == nil {
		logEntry.Error("url field not found or incorrect type")
	} else {
		url := fmt.Sprintf("%v", config["url"])
		c.URL = url
	}

	if config["body"] == nil {
		logEntry.Error("body field not found or incorrect type")
	} else {
		body := fmt.Sprintf("%v", config["body"])
		c.Body = body
	}
}

func NewHTTPOutput() interface{} {
	o := &HTTPOutput{
		nil,
		make(chan []byte),
		&http.Client{},
	}

	httpOutputLogger = logrus.WithField("output", o)

	return o
}

func (o *HTTPOutput) SetConfig(c model.Config) bool {
	if c.GetType() != HTTPOutputType || c.GetName() == "" {
		return false
	}

	if c := c.GetUserConfig(); c == nil {
		return false
	}

	o.config = c
	return true
}

func (o *HTTPOutput) GetConfig() model.Config {
	return o.config
}

func (o *HTTPOutput) Execute() bool {
	httpOutputConfig := o.config.(*HTTPOutputConfig)
	request, err := http.NewRequest(httpOutputConfig.Method, httpOutputConfig.URL, strings.NewReader(httpOutputConfig.Body))

	if err != nil {
		httpOutputLogger.WithField("error", err).Error("Malformed HTTP request")
		return false
	}

	res, err := o.httpClient.Do(request)
	if err != nil {
		httpOutputLogger.WithFields(logrus.Fields{"error": err}).Error("Output error")
	}

	httpOutputLogger.WithFields(logrus.Fields{"response": res}).Debug("Got HTTP output response")

	return err == nil
}
