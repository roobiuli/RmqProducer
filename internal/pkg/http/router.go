package http

import (
	"github.com/roobiuli/RmqProducer/internal/pkg/rabbitmq"
	"github.com/roobiuli/RmqProducer/internal/user"
	"net/http"
)

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{http.NewServeMux()}
}

func (r *Router) RegisterUser(rb *rabbitmq.RabbitMQ)  {
	create := user.NewCreate(rb)
	r.HandleFunc("/users/create", create.Handler)
}