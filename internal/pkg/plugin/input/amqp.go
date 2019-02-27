package input

import (
	"fmt"
	"github.com/CultOfAgileLimericks/conductor/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var amqpInputLogger *logrus.Entry

const AMQPInputType = "amqp"

type AMQPInput struct {
	config model.Config
	channel chan<-model.Input
	stop chan bool
}

type AMQPInputConfig struct {
	name string
	ConnectionString string
	ExchangeName string
	Topic string
	QueueName string
}

func (*AMQPInputConfig) GetType() string {
	return AMQPInputType
}

func (config *AMQPInputConfig) GetName() string {
	return config.name
}

func (config *AMQPInputConfig) SetName(name string) {
	config.name = name
}

func (config *AMQPInputConfig) GetUserConfig() map[string]interface{} {
	if config.ConnectionString == "" || config.ExchangeName == "" || config.Topic == "" {
		return nil
	}

	userConfig := make(map[string]interface{})
	userConfig["connectionString"] = config.ConnectionString
	userConfig["exchangeName"] = config.ExchangeName
	userConfig["topic"] = config.Topic
	userConfig["queueName"] = config.QueueName

	return userConfig
}

func (config *AMQPInputConfig) SetUserConfig(c map[string]interface{}) {
	logEntry := logrus.WithField("config", config)

	if c["connectionString"] == nil {
		logEntry.Error("connectionString field not found or incorrect type")
	} else {
		config.ConnectionString = fmt.Sprintf("%v", c["connectionString"])
	}

	if c["exchangeName"] == nil {
		logEntry.Error("exchangeName field not found or incorrect type")
	} else {
		config.ExchangeName = fmt.Sprintf("%v", c["exchangeName"])
	}

	if c["topic"] == nil {
		logEntry.Error("topic field not found or incorrect type")
	} else {
		config.Topic = fmt.Sprintf("%v", c["topic"])
	}

	if c["queueName"] == nil {
		logEntry.Info("queueName field not found or incorrect type")
	} else {
		config.QueueName = fmt.Sprintf("%v", c["queueName"])
	}
}

func NewAMQPInput() interface{} {
	amqpInput := &AMQPInput{
		nil,
		nil,
		make(chan bool),
	}

	amqpInputLogger = logrus.WithField("input", amqpInput)

	return amqpInput
}

func (i *AMQPInput) SetConfig(c model.Config) bool {
	if c.GetName() == "" || c.GetType() != AMQPInputType {
		return false
	}

	if userConfig := c.GetUserConfig(); userConfig == nil {
		return false
	}

	i.config = c
	return true
}

func (i *AMQPInput) GetConfig() model.Config {
	return i.config
}

func (i *AMQPInput) SetInputChannel(c chan<- model.Input) {
	i.channel = c
}

func (i *AMQPInput) Listen() {
	conn, err := amqp.Dial(i.GetConfig().(*AMQPInputConfig).ConnectionString)
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to open AMQP connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to get AMQP channel")
	}

	err = ch.ExchangeDeclare(
		i.GetConfig().(*AMQPInputConfig).ExchangeName,
		i.GetConfig().(*AMQPInputConfig).Topic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to declare exchange")
	}

	q, err := ch.QueueDeclare(
		i.GetConfig().(*AMQPInputConfig).QueueName,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to declare queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to consume messages")
	}

	RangeLoop:
	for {
		select {
		case <- msgs:
			i.channel <- i
		case <- i.stop:
			break RangeLoop
		}
	}

	err = conn.Close()
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to close AMQP connection")
	}

	err = ch.Close()
	if err != nil {
		amqpInputLogger.WithField("err", err).Error("Failed to close AMQP channel")
	}
}

func (i *AMQPInput) Stop() {
	i.stop <- true
}
