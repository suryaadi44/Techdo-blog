package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

type Server struct {
	Address string
	Router  *mux.Router
}

func (s *Server) Run() {
	go func() {
		if err := http.ListenAndServe(s.Address, s.Router); err != nil {
			log.Printf("[SERVER] error starting server at %s: %v\n", s.Address, err)
		}
	}()

	log.Printf("[SERVER] server started, listening to %s\n", s.Address)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Printf("[SERVER] server stopped")
}

func NewServer(address string, router *mux.Router) Server {
	return Server{
		Address: address,
		Router:  router,
	}
}
