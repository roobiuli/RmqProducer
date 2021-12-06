package config

import (
	"github.com/roobiuli/RmqProducer/internal/pkg/rabbitmq"
	"github.com/roobiuli/RmqProducer/internal/user"
	"time"
)

type Config struct {
	HTTPAddress string
	RabbitMQ rabbitmq.Config
	UserAMPQ user.AMQPConfig
}


func New() Config {
	var cnf Config

	cnf.HTTPAddress = ":8080"
	cnf.RabbitMQ.Schema = "amqp"
	cnf.RabbitMQ.Username = "PatatoUser"
	cnf.RabbitMQ.Password = "123123"
	cnf.RabbitMQ.Host = "0.0.0.0"
	cnf.RabbitMQ.Port = "5672"
	cnf.RabbitMQ.Vhost = "my_app"
	cnf.RabbitMQ.ConnectionName = "MY_APP"
	cnf.RabbitMQ.ChannelNotifyTimout = 100 * time.Millisecond
	cnf.RabbitMQ.Reconnect.Interval = 500 * time.Millisecond
	cnf.RabbitMQ.Reconnect.MaxAttempt = 100

	cnf.UserAMPQ.Create.ExchangeName = "user"
	cnf.UserAMPQ.Create.EXchangeType = "direct"
	cnf.UserAMPQ.Create.RoutingKey = "create"
	cnf.UserAMPQ.Create.QueueName = "user_create"
	return cnf
}