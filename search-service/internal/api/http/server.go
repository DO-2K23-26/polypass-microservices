package http

import (
	httpController "github.com/DO-2K23-26/polypass-microservices/search-service/controller/http"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server *gin.Engine
	port   string
}

func NewServer(healthController *httpController.HealthController, port string) *Server {
	server := gin.Default()
	server.GET("/health", healthController.CheckHealth)
	return &Server{
		server: server,
		port:   port,
	}
}

func (s *Server) Start() error {
	return s.server.Run(":" + s.port)
}
