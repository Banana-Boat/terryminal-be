package http

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) conversation(ctx *gin.Context) {

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	messagechan := make(chan string)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			now := time.Now().Format("2006-01-02 15:04:05")
			currentTime := fmt.Sprintf("The Current Time Is %v", now)

			// Send current time to clients message channel
			fmt.Println(currentTime)
			messagechan <- currentTime
		}
	}()

	defer func() {
		close(messagechan)
		messagechan = nil
		log.Print("client connection is closed")
	}()

	ctx.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-messagechan; ok {
			ctx.SSEvent("message", msg)
			return true
		}
		return false
	})
}
