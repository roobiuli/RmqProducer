package user

import (
	"errors"
	"github.com/roobiuli/RmqProducer/internal/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

type Create struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func NewCreate(r *rabbitmq.RabbitMQ) Create {
	return Create{
		rabbitmq: r,
	}
}

func (c *Create) Handler(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("ID")

		if err := c.Publish(id); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Created", id)
	return
}

func (c *Create) Publish(message string) error  {
	channel, err := c.rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	if err := channel.Confirm(false) ; err != nil {
		log.Println(err)
	}

	if err := channel.Publish(
		"user",
		"create",
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId: "A-UNIQUE-ID",
			ContentType: "text/plain",
			Body: []byte(message),
		},
		); err != nil {
		return err
	}

	select {
		case ntf := <-channel.NotifyPublish(make(chan amqp.Confirmation, 1)):
			if !ntf.Ack {
				return errors.New("failed to deliver message to exchange/queue")
			}
		case <-channel.NotifyReturn(make(chan amqp.Return)):
			return errors.New("failed to deliver message to exchange/queue")
		case <-time.After(c.rabbitmq.ChannelNotifyTimeout):
			log.Println("message delivery confirmation to exchange/queue timed out")
		}

		return nil
	}

