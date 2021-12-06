package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type Config struct {
	Schema string
	Username string
	Password string
	Port  string
	Vhost string
	Host string
	ConnectionName string
	ChannelNotifyTimout time.Duration
	Reconnect struct{
		Interval time.Duration
		MaxAttempt int
	}
}

type RabbitMQ struct {
	mux sync.RWMutex
	config Config
	dialConfig amqp.Config
	ChannelNotifyTimeout time.Duration
	connection *amqp.Connection
}


func New(conf Config) *RabbitMQ {
	return &RabbitMQ{
		config:               conf,
		dialConfig:           amqp.Config{Properties: amqp.Table{"connection_name": conf.ConnectionName}},
		ChannelNotifyTimeout: conf.ChannelNotifyTimout,
	}
}

func(r *RabbitMQ) Connect() error {
	con, err := amqp.DialConfig(fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		r.config.Schema,
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port,
		r.config.Vhost,
		), r.dialConfig)

	if err != nil {
		return err
	}

	r.connection = con

	go r.Reconnect()
	return nil
}

func (r *RabbitMQ) Channel() (*amqp.Channel, error)  {
	if r.connection == nil {
		if err := r.Connect() ; err != nil {
			return nil, errors.New("Connection is not open")
		}
	}
	channel, err := r.connection.Channel()

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *RabbitMQ) Connection() *amqp.Connection  {
	return r.connection
}


func (r *RabbitMQ) ShutDown() error  {
	if r.connection != nil {
		r.connection.Close()
	}
	return nil
}


func (r *RabbitMQ) Reconnect() {

	conErr := <-r.connection.NotifyClose(make(chan *amqp.Error))
	if conErr != nil {
		log.Println("CRITICIAL: Connection dropped, reconnecting")

		var err error
		for i := 1 ; i <= r.config.Reconnect.MaxAttempt; i++ {
			r.mux.Lock()
			r.connection, err = amqp.DialConfig(fmt.Sprintf(
				"%s://%s:%s@%s:%s/%s",
				r.config.Schema,
				r.config.Username,
				r.config.Password,
				r.config.Host,
				r.config.Port,
				r.config.Vhost,
			), r.dialConfig)
			r.mux.RUnlock()
		}

		if err != nil {
			log.Println("INFO: Reconnected")
		}

	}

}

