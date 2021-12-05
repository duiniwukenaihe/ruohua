package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	self_jwt "github.com/chenhqchn/ruohua/server/utils/jwt"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authValue := ctx.Request.Header.Get("Authorization")
		if authValue == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 401, "msg": "The request you must be made requires authentication"})
			ctx.Abort()
			return
		}

		claims, err := self_jwt.ParseToken(authValue)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 401, "msg": fmt.Sprintf("%s", err)})
			ctx.Abort()
			return
		}

		ctx.Set("userClaims", claims)
		ctx.Next()
	}
}