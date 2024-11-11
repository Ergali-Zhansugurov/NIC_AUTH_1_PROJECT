package http

import (
	"awesomeProject4/user-auth-service/internal/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	Router *gin.Engine
	Config *config.Config
}

// NewServer инициализирует сервер с настройками
func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	return &Server{
		Router: router,
		Config: cfg,
	}
}

// Run запускает сервер
func (s *Server) Run() error {
	port := s.Config.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := s.Router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Could not start server: %v", err)
		return err
	}
	return nil
}
