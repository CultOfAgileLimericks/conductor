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
	config model.Config
	channel chan <-[]byte
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

	userConfig["Method"] = c.Method
	userConfig["URL"] = c.URL
	userConfig["Body"] = c.Body

	return userConfig
}


func (c *HTTPOutputConfig) SetUserConfig(config map[string]interface{}) {
	method, ok := config["Method"].(string)
	logEntry := logrus.WithField("config", c)
	if !ok {
		logEntry.Error("Method field not found or incorrect type")
	}
	c.Method = method

	url, ok := config["URL"].(string)
	if !ok {
		logEntry.Error("URL field not found or incorrect type")
	}
	c.URL = url

	body, ok := config["Body"].(string)
	if !ok {
		logEntry.Error("Body field not found or incorrect type")
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
