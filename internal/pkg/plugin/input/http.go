package input

import (
	"log"
	"net/http"
	"strings"
)

type HTTPInput struct {
	Addr string
	channel chan <-[]byte
}

func NewHTTPInput(addr string) *HTTPInput {
	return &HTTPInput{
		addr,
		nil,
	}
}

func (input *HTTPInput) SetInputChannel(c chan <-[]byte) {
	input.channel = c
}

func (input *HTTPInput) Listen() {
	msg := make([]byte, 3)
	_, _ = strings.NewReader("asd").Read(msg)
	// TODO: Send something useful to this channel
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- msg
	})
	log.Fatal(http.ListenAndServe(input.Addr, nil))
}
