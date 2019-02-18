package output

import "testing"

func TestNewHTTPOutput(t *testing.T) {
	o := NewHTTPOutput()

	if o.GetConfig() != nil {
		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputName(t *testing.T) {
	c := &HTTPOutputConfig{}
	c.SetName("test output")

	if c.GetName() != "test output" {
		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputUserConfig_Correct(t *testing.T) {
	c := &HTTPOutputConfig{}
	userConfig := make(map[string]interface{})

	userConfig["Method"] = "GET"
	userConfig["URL"] = "http://test.com"
	userConfig["Body"] = "someBODY once told me..."

	c.SetUserConfig(userConfig)

	if c.Method != userConfig["Method"] ||
		c.URL != userConfig["URL"] ||
		c.Body != userConfig["Body"] {

		t.Fail()
	}
}

func TestHTTPOutputConfig_SetOutputUserConfig_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{}
	userConfig := make(map[string]interface{})

	c.SetUserConfig(userConfig)

	if c.Method != "" || c.GetName() != "" || c.Body != "" {
		t.Fail()
	}
}

func TestHTTPOutput_UseConfig_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{}
	h := NewHTTPOutput()

	ok := h.SetConfig(c)

	if ok || h.GetConfig() != nil {
		t.Fail()
	}

	c.SetName("test output")
	ok = h.SetConfig(c)

	if ok || h.GetConfig() != nil {
		t.Fail()
	}
}

func TestHTTPOutput_UseConfig_Correct(t *testing.T) {
	c := &HTTPOutputConfig{
		Method: "GET",
		URL:    "https://google.com",
		Body:   "someBODY once told me...",
	}
	c.SetName("test config")

	h := NewHTTPOutput()

	ok := h.SetConfig(c)

	if !ok || h.GetConfig() != c {
		t.Fail()
	}
}

func TestHTTPOutput_Execute_Correct(t *testing.T) {
	c := &HTTPOutputConfig{
		Method: "GET",
		URL:    "https://google.com",
		Body:   "",
	}
	c.SetName("test config")

	h := NewHTTPOutput()
	h.SetConfig(c)

	ok := h.Execute()

	if !ok {
		t.Fail()
	}
}

func TestHTTPOutput_Execute_Incorrect(t *testing.T) {
	c := &HTTPOutputConfig{
		Method: "GET",
		URL:    "lolololololololololololololol",
		Body:   "",
	}
	c.SetName("test config")

	h := NewHTTPOutput()
	h.SetConfig(c)

	ok := h.Execute()

	if ok {
		t.Fail()
	}

	c.Method = "THIS IS AN UNKNOWN METHOD"
	h = NewHTTPOutput()
	h.SetConfig(c)

	ok = h.Execute()

	if ok {
		t.Fail()
	}
}
