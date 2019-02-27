package input

import (
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"testing"
)

func TestNewAMQPInput(t *testing.T) {
	c := NewAMQPInput().(*AMQPInput)
	if c.GetConfig() != nil {
		t.Fail()
	}
}

func TestAMQPInputConfig_SetInputName(t *testing.T) {
	config := &AMQPInputConfig{}
	config.SetName("test name")

	if config.GetName() != "test name" {
		t.Fail()
	}
}

func TestAMQPInputConfig_SetInputUserConfig_Correct(t *testing.T) {
	config := &AMQPInputConfig{}
	userConfig := make(map[string]interface{})
	userConfig["connectionString"] = "amqp://127.0.0.1:5672"
	userConfig["exchangeName"] = "test exchange"
	userConfig["topic"] = "test topic"
	userConfig["queueName"] = "test queue"
	config.SetUserConfig(userConfig)

	if config.GetUserConfig()["connectionString"] != userConfig["connectionString"] ||
		config.GetUserConfig()["exchangeName"] != userConfig["exchangeName"] ||
		config.GetUserConfig()["topic"] != userConfig["topic"] ||
		config.GetUserConfig()["queueName"] != userConfig["queueName"] {
		t.Fail()
	}

	if config.ConnectionString != userConfig["schedule"] ||
		config.ExchangeName != userConfig["exchangeName"] ||
		config.Topic != userConfig["topic"] ||
		config.QueueName != userConfig["queueName"] {
		t.Fail()
	}
}

func TestAMQPInputConfig_SetInputUserConfig_Incorrect(t *testing.T) {
	config := &AMQPInputConfig{}
	userConfig := make(map[string]interface{})
	config.SetUserConfig(userConfig)

	if config.ConnectionString != ""  ||
		config.ExchangeName != "" ||
		config.Topic != "" ||
		config.QueueName != "" {
		t.Fail()
	}

	if config.GetUserConfig()["connectionString"] != nil ||
		config.GetUserConfig()["exchangeName"] != nil ||
		config.GetUserConfig()["topic"] != nil ||
		config.GetUserConfig()["queueName"] != nil {
		t.Fail()
	}
}

func TestAMQPInput_UseConfig_Incorrect(t *testing.T) {
	c := NewAMQPInput().(*AMQPInput)
	config := &AMQPInputConfig{}

	ok := c.SetConfig(config)
	if ok || c.config == config {
		t.Fail()
	}

	config.SetName("TestAMQPInput")

	ok = c.SetConfig(config)
	if ok || c.GetConfig() == config {
		t.Fail()
	}
}

func TestAMQPInput_UseConfig_Correct(t *testing.T) {
	c := NewAMQPInput().(*AMQPInput)
	config := &AMQPInputConfig{
		ConnectionString: "amqp://localhost:5672",
		ExchangeName: "test exchange",
		Topic: "test topic",
		QueueName: "test queue",
	}
	config.SetName("TestAMQPInput")

	ok := c.SetConfig(config)
	if !ok || c.GetConfig() != config {
		t.Fail()
	}
}

func TestAMQPInput_SetInputChannel(t *testing.T) {
	c := NewAMQPInput().(*AMQPInput)

	channel := make(chan model.Input)
	c.SetInputChannel(channel)

	if c.channel != channel {
		t.Fail()
	}
}
//
//func TestAMQPInput_Listen_Correct(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	defer ctrl.Finish()
//
//	c := NewAMQPInput().(*AMQPInput)
//	config := &AMQPInputConfig{
//		ConnectionString: "amqp://localhost:5672",
//		ExchangeName: "test exchange",
//		Topic: "test topic",
//		QueueName: "test queue",
//	}
//
//	c.SetConfig(config)
//	channel := make(chan model.Input)
//
//	c.SetInputChannel(channel)
//
//	go c.Listen()
//
//	defer c.Stop()
//
//	select {
//	case i := <-channel:
//		if i != c {
//			t.Fail()
//		}
//	case <-time.After(2 * time.Second):
//		t.Fail()
//	}
//}

