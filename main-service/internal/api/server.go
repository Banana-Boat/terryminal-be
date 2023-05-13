package api

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker *util.TokenMaker
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := util.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) Start() error {
	if err := server.router.Run(
		fmt.Sprintf("%s:%s", server.config.MainServerHost, server.config.MainServerPort),
	); err != nil {
		return err
	}

	return nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user/login", server.login)
	router.POST("/user/register", server.register)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.GET("/user/listUsers", server.listUsers)

	// authRouter.POST("/chatbot/chat", server.chat)
	router.POST("/chatbot/chat", server.chat)

	server.router = router
}

func wrapResponse(flag bool, msg string, data gin.H) gin.H {
	var _msg string
	if msg == "" {
		_msg = "OK"
	} else {
		_msg = msg
	}

	return gin.H{
		"flag": flag,
		"msg":  _msg,
		"data": data,
	}
}
