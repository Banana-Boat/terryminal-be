package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var err = errors.New("invalid authorization header format")

func authMiddleware(tokenMaker *TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("authorization")

		if len(authHeader) == 0 {
			log.Error().Err(err).Msg("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			log.Error().Err(err).Msg("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != "bearer" {
			log.Error().Err(err).Msg("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			log.Error().Err(err).Msg("failed to verify token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, wrapResponse(false, err.Error(), nil))
			return
		}

		ctx.Set("token_payload", payload)
		ctx.Next()
	}
}
