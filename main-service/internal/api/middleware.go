package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func authMiddleware(tokenMaker *TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("authorization")

		if len(authHeader) == 0 {
			log.Error().Msg("authorization header is not provided")
			ctx.JSON(http.StatusBadRequest, wrapResponse(false, "请求头错误", nil))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			log.Error().Msg("invalid authorization header format")
			ctx.JSON(http.StatusBadRequest, wrapResponse(false, "请求头错误", nil))
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != "bearer" {
			log.Error().Msg("invalid authorization header format")
			ctx.JSON(http.StatusBadRequest, wrapResponse(false, "请求头错误", nil))
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			log.Error().Msg("failed to verify token")
			ctx.JSON(http.StatusUnauthorized, wrapResponse(false, "鉴权失败", nil))
			return
		}

		ctx.Set("token_payload", payload)
		ctx.Next()
	}
}
