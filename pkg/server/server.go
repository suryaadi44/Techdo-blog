package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	Router *mux.Router
}

func (s *Server) Run() {
	host := fmt.Sprintf(":%s", s.Port)

	go func() {
		if err := http.ListenAndServe(host, s.Router); err != nil {
			log.Printf("[SERVER] error starting server at %s: %v\n", host, err)
		}
	}()

	log.Printf("[SERVER] server started, listening to %s\n", host)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Printf("[SERVER] server stopped")
}

func NewServer(port string, router *mux.Router) Server {
	return Server{
		Port:   port,
		Router: router,
	}
}
