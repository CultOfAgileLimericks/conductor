package output

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var httpOutputLogger *logrus.Entry

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
	return "http"
}

func (c *HTTPOutputConfig) OutputName() string {
	return c.Name
}

func (c *HTTPOutputConfig) SetOutputName(name string) {
	c.Name = name
}

func (c *HTTPOutputConfig) OutputUserConfig() map[string]interface{} {
	userConfig := make(map[string]interface{})

	userConfig["method"] = c.Method
	userConfig["url"] = c.URL
	userConfig["body"] = c.Body

	return userConfig
}


func (c *HTTPOutputConfig) SetOutputUserConfig(config map[string]interface{}) {
	method, ok := config["method"].(string)
	if !ok {
		httpOutputLogger.Fatal("method field not found or incorrect type")
	}
	c.Method = method

	url, ok := config["url"].(string)
	if !ok {
		httpOutputLogger.Fatal("url field not found or incorrect type")
	}
	c.URL = url

	body, ok := config["body"].(string)
	if !ok {
		httpOutputLogger.Fatal("body field not found or incorrect type")
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

func (o *HTTPOutput) UseConfig(c model.OutputConfig) {
	o.Config = c
}

func (o *HTTPOutput) Execute() bool {
	httpOutputConfig := o.Config.(*HTTPOutputConfig)
	request, err := http.NewRequest(httpOutputConfig.Method, httpOutputConfig.URL, strings.NewReader(httpOutputConfig.Body))

	if err != nil {
		httpOutputLogger.WithField("error", err).Error("Malformed HTTP request")
	}

	res, err := o.httpClient.Do(request)
	if err != nil {
		httpOutputLogger.WithFields(logrus.Fields{"error": err}).Error("Output error")
	}

	httpOutputLogger.WithFields(logrus.Fields{"response": res}).Debug("Got HTTP output response")

	return err == nil
}
