package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"log"
	"net/http"
)

type HTTPInput struct {
	Addr string
	channel chan <-model.Input
}

func NewHTTPInput(addr string) *HTTPInput {
	return &HTTPInput{
		addr,
		nil,
	}
}

func (input *HTTPInput) SetInputChannel(c chan <-model.Input) {
	input.channel = c
}

func (input *HTTPInput) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		input.channel <- input
	})
	log.Fatal(http.ListenAndServe(input.Addr, nil))
}
