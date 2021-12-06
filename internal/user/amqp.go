package user

import (
	"errors"
	"github.com/roobiuli/RmqProducer/internal/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

type AMQPConfig struct {
	Create struct{
		ExchangeName string
		EXchangeType string
		RoutingKey string
		QueueName string
	}
}


type AMPQ struct {
	config AMQPConfig
	rabbitmq *rabbitmq.RabbitMQ
}


func NewAMPQ(config AMQPConfig, rm *rabbitmq.RabbitMQ) AMPQ {
	return AMPQ{
		config: config,
		rabbitmq: rm,
	}
}

func (a *AMPQ) Setup() error {
	channel, err := a.rabbitmq.Channel()
	if err != nil {
		return errors.Unwrap(err)
	}
	defer channel.Close()
	if err := a.DeclareCreate(channel); err != nil {
		return err
	}
	return nil
}

func (a *AMPQ) DeclareCreate(channel *amqp.Channel) error {
	if err := channel.ExchangeDeclare(
		a.config.Create.ExchangeName,
		a.config.Create.EXchangeType,
		true,
		false,
		false,
		false,
		nil,
		); err != nil {
		return err
	}
	if err := channel.QueueBind(
		a.config.Create.QueueName,
		a.config.Create.RoutingKey,
		a.config.Create.ExchangeName,
		false,
		nil,
		); err != nil {
		return err
	}
	return nil
}



