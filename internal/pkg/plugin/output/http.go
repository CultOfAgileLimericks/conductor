package output

type HTTPOutput struct {
	channel chan <-[]byte
}

func NewHTTPOutput() *HTTPOutput {
	return &HTTPOutput{}
}

func (o *HTTPOutput) Execute() bool {
	return true
}
