/*
Package auth provides Auth middlewares.
*/

package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexanderbkl/vidre-back/internal/api"
	"github.com/alexanderbkl/vidre-back/internal/constant"
	"github.com/alexanderbkl/vidre-back/pkg/token"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(constant.AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != constant.AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		ctx.Set(constant.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
