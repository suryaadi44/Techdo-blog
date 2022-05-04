package server

import (
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Address string
	Engine  *gin.Engine
}

func (s *Server) Run() {
	go func() {
		if err := s.Engine.Run(s.Address); err != nil {
			log.Printf("[SERVER] error starting server at %s: %v\n", s.Address, err)
		}
	}()

	log.Printf("[SERVER] server started, listening to %s\n", s.Address)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Printf("[SERVER] server stopped")
}

func NewServer(address string, engine *gin.Engine) Server {
	return Server{
		Address: address,
		Engine:  engine,
	}
}
