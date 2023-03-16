package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	DisplayName string
	UsrId       int64
	Role        string
}

func Jwt() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if ctx.Request.Method == http.MethodGet {
			ctx.Next()
			return
		}

		user := ctx.GetHeader("Authorization")
		if len(user) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		ctx.Set("user_id", int64(2))
		ctx.Next()
	}

}
