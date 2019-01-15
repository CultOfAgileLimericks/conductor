package output

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger *logrus.Entry

type HTTPOutput struct {
	HTTPRequest *http.Request
	channel chan <-[]byte
	httpClient *http.Client
}

func NewHTTPOutput(r *http.Request) *HTTPOutput {
	o := &HTTPOutput{
		r,
		make(chan []byte),
		&http.Client{},

	}

	logger = logrus.WithField("output", o)

	return o
}

func (o *HTTPOutput) Execute() bool {
	res, err := o.httpClient.Do(o.HTTPRequest)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("Output error")
	}

	logger.WithFields(logrus.Fields{"response": res}).Debug("Got HTTP output response")

	return err == nil
}
