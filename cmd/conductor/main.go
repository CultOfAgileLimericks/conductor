package main

import (
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/input"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/plugin/output"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	t := model.NewTask()
	httpInput := input.NewHTTPInput(":8080")

	r, _ := http.NewRequest("GET", "https://google.com", nil)

	httpOutput := output.NewHTTPOutput(r)

	t.RegisterInput(httpInput)
	t.RegisterOutput(httpOutput)

	t.Run()

	fmt.Printf("%+v", t)
}
