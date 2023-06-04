package api

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/Banana-Boat/terryminal/main-service/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type Server struct {
	config          util.Config
	store           *db.Store
	tokenMaker      *TokenMaker
	taskDistributor *worker.TaskDistributor
	router          *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	}
	taskDistributor := worker.NewTaskDistributor(redisOpt)

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
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

	router.POST("/user/login", server.loginHandle)
	router.POST("/user/register", server.registerHandle)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.PATCH("/user/updateInfo", server.updateInfoHandle)

	router.POST("/chatbot/chat", server.chatHandle)
	router.GET("/terminal/ws", server.terminalWSHandle)

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
