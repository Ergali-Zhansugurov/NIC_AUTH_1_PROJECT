package http

import (
	"awesomeProject4/user-auth-service/internal/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
func (s *Server) Run() {
	port := s.Config.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := s.Router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// SetupRoutes задает маршруты (временно добавим тестовый маршрут)
func (s *Server) SetupRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
