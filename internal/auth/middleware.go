package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

type User struct {
	DisplayName string `json:"display_name"`
	UserId      int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJwtSecret() {

	jwtSecret = os.Getenv("AUTH_SECRET")
	if len(jwtSecret) <= 0 {
		log.Fatal("variable 'AUTH_SECRET' in env file is empty")
	}

}

func Jwt() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		if ctx.Request.Method == http.MethodGet {
			ctx.Next()
			return
		}

		tokenString := ctx.GetHeader("Authorization")
		header := strings.Split(tokenString, " ")

		if len(tokenString) == 0 || len(header) != 2 || header[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": http.StatusText(http.StatusUnauthorized),
				"data":    nil,
			})
			return
		}

		var token *jwt.Token
		var err error

		// Просто заглушка для тестов :)
		if gin.Mode() != gin.ReleaseMode {
			ctx.Set("user_id", int64(2))
			ctx.Next()
			return
		}

		errorChecking := map[string]any{
			"status": "error",
			"data":   nil,
		}

		if token, err = jwt.ParseWithClaims(tokenString, &User{}, parseJwtMethod); err != nil || !token.Valid {
			errorChecking["message"] = "Invalid token"
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H(errorChecking))
			return
		}

		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			errorChecking["message"] = "Bad token"
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			errorChecking["message"] = "Invalid signature"
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			errorChecking["message"] = "Token is expired"
		default:
			if claims, ok := token.Claims.(*User); ok {
				ctx.Set("user_id", claims.UserId)
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H(errorChecking))

	}

}

func parseJwtMethod(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("invalid signing method")
	}
	return []byte(jwtSecret), nil

}
