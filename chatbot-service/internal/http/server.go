package http

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/chatbot-service/internal/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	router *gin.Engine
}

func NewServer(config util.Config) *Server {
	server := &Server{
		config: config,
	}

	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/chat", server.chat)

	server.router = router
}

func (server *Server) Start() error {
	if err := server.router.Run(
		fmt.Sprintf("%s:%s", server.config.ChatbotHttpServerHost, server.config.ChatbotHttpServerPort),
	); err != nil {
		return err
	}

	return nil
}
