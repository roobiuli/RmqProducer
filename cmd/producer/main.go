package main

import (
	app "github.com/roobiuli/RmqProducer/internal/app"
	"github.com/roobiuli/RmqProducer/internal/config"
	"github.com/roobiuli/RmqProducer/internal/pkg/http"
	"github.com/roobiuli/RmqProducer/internal/pkg/rabbitmq"
	"github.com/roobiuli/RmqProducer/internal/user"
	"log"

)

func main() {
	// Config

	cnf := config.New()

	// RBTMq conn

	rb := rabbitmq.New(cnf.RabbitMQ)

	if err := rb.Connect() ; err != nil {
		log.Fatalln(err)
	}

	defer rb.ShutDown()

	// AMQP Setup
	userMPQP := user.NewAMPQ(cnf.UserAMPQ, rb)

	if err := userMPQP.Setup() ; err != nil {
		log.Fatalln(err)
	}

	// setup Router
	router := http.NewRouter()

	router.RegisterUser(rb)

	// SRV

	srv := http.NewServer(cnf.HTTPAddress, router)

	// app RUN
	if err := app.NewApp(srv).Run() ; err != nil {
		log.Fatalln(err)
	}
}