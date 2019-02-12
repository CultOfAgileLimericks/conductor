package output

import "testing"

func TestNewHTTPOutput(t *testing.T) {
	o := NewHTTPOutput()

	if o.Config != nil {
		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputName(t *testing.T) {
	c := &HTTPOutputConfig{}
	c.SetOutputName("test output")

	if c.Name != "test output" {
		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputUserConfig_Correct(t *testing.T) {
	c := &HTTPOutputConfig{}
	userConfig := make(map[string]interface{})

	userConfig["method"] = "GET"
	userConfig["url"] = "http://test.com"
	userConfig["body"] = "someBODY once told me..."

	c.SetOutputUserConfig(userConfig)

	if c.Method != userConfig["method"] ||
		c.URL != userConfig["url"] ||
		c.Body != userConfig["body"] {

		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputUserConfig_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{}
	userConfig := make(map[string]interface{})

	c.SetOutputUserConfig(userConfig)

	if c.Method != "" || c.Name != "" || c.Body != "" {
		t.Fail()
	}
}

func TestHTTPOutput_UseConfig_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{}
	h := NewHTTPOutput()

	ok := h.UseConfig(c)

	if ok || h.Config != nil {
		t.Fail()
	}

	c.SetOutputName("test output")
	ok = h.UseConfig(c)

	if ok || h.Config != nil {
		t.Fail()
	}
}

func TestHTTPOutput_UseConfig_Correct(t *testing.T) {
	c := &HTTPOutputConfig{
		Name: "test config",
		Method: "GET",
		URL: "https://google.com",
		Body: "someBODY once told me...",
	}

	h := NewHTTPOutput()

	ok := h.UseConfig(c)

	if !ok || h.Config != c {
		t.Fail()
	}
}

func TestHTTPOutput_Execute_Correct(t *testing.T) {
	c := &HTTPOutputConfig{
		Name: "test config",
		Method: "GET",
		URL: "https://google.com",
		Body: "",
	}

	h := NewHTTPOutput()
	h.UseConfig(c)

	ok := h.Execute()

	if !ok {
		t.Fail()
	}
}

func TestHTTPOutput_Execute_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{
		Name: "test config",
		Method: "GET",
		URL: "lolololololololololololololol",
		Body: "",
	}

	h := NewHTTPOutput()
	h.UseConfig(c)

	ok := h.Execute()

	if ok {
		t.Fail()
	}

	c.Method = "THIS IS AN UNKNOWN METHOD"
	h = NewHTTPOutput()
	h.UseConfig(c)

	ok = h.Execute()

	if ok {
		t.Fail()
	}
}
