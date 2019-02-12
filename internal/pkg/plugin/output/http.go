package output

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var httpOutputLogger *logrus.Entry
const HTTPOutputType = "http"

type HTTPOutput struct {
	Config model.OutputConfig
	channel chan <-[]byte
	httpClient *http.Client
}

type HTTPOutputConfig struct {
	Name string
	Method string
	URL string
	Body string
}

func (c *HTTPOutputConfig) OutputType() string {
	return HTTPOutputType
}

func (c *HTTPOutputConfig) OutputName() string {
	return c.Name
}

func (c *HTTPOutputConfig) SetOutputName(name string) {
	c.Name = name
}

func (c *HTTPOutputConfig) OutputUserConfig() map[string]interface{} {
	if c.URL == "" || c.Method == "" {
		return nil
	}

	userConfig := make(map[string]interface{})

	userConfig["method"] = c.Method
	userConfig["url"] = c.URL
	userConfig["body"] = c.Body

	return userConfig
}


func (c *HTTPOutputConfig) SetOutputUserConfig(config map[string]interface{}) {
	method, ok := config["method"].(string)
	logEntry := logrus.WithField("config", c)
	if !ok {
		logEntry.Error("method field not found or incorrect type")
	}
	c.Method = method

	url, ok := config["url"].(string)
	if !ok {
		logEntry.Error("url field not found or incorrect type")
	}
	c.URL = url

	body, ok := config["body"].(string)
	if !ok {
		logEntry.Error("body field not found or incorrect type")
	}
	c.Body = body
}

func NewHTTPOutput() *HTTPOutput {
	o := &HTTPOutput{
		nil,
		make(chan []byte),
		&http.Client{},
	}

	httpOutputLogger = logrus.WithField("output", o)

	return o
}

func (o *HTTPOutput) UseConfig(c model.OutputConfig) bool {
	if c.OutputType() != HTTPOutputType || c.OutputName() == "" {
		return false
	}

	if c := c.OutputUserConfig(); c == nil {
		return false
	}

	o.Config = c
	return true
}

func (o *HTTPOutput) Execute() bool {
	httpOutputConfig := o.Config.(*HTTPOutputConfig)
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
