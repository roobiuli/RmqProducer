package app

import (
	"net/http"
)


type App struct {
	server *http.Server
}

func NewApp(serv *http.Server) App {
	return App{
		server: serv,
		}
}


func(a App) Run() error {
	return a.server.ListenAndServe()
}

