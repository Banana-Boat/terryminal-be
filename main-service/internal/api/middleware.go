package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Banana-Boat/go-micro-template/main-service/internal/util"
	"github.com/gin-gonic/gin"
)

var err = errors.New("invalid authorization header format")

func authMiddleware(tokenMaker *util.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("authorization")

		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		ctx.Set("auth_payload", payload)
		ctx.Next()
	}
}
