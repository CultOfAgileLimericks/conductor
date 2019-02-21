package output

import "testing"

func TestNewHTTPOutput(t *testing.T) {
	o := NewHTTPOutput().(*HTTPOutput)

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

	userConfig["method"] = "GET"
	userConfig["url"] = "http://test.com"
	userConfig["body"] = "someBODY once told me..."

	c.SetUserConfig(userConfig)

	if c.Method != userConfig["method"] ||
		c.URL != userConfig["url"] ||
		c.Body != userConfig["body"] {

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
	h := NewHTTPOutput().(*HTTPOutput)

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

	h := NewHTTPOutput().(*HTTPOutput)

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

	h := NewHTTPOutput().(*HTTPOutput)
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

	h := NewHTTPOutput().(*HTTPOutput)
	h.SetConfig(c)

	ok := h.Execute()

	if ok {
		t.Fail()
	}

	c.Method = "THIS IS AN UNKNOWN METHOD"
	h = NewHTTPOutput().(*HTTPOutput)
	h.SetConfig(c)

	ok = h.Execute()

	if ok {
		t.Fail()
	}
}
