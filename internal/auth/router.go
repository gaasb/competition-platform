package auth

import "github.com/gin-gonic/gin"

type AuthRouter struct{}

func (r *AuthRouter) Setup(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", handleLogin)
		auth.POST("/logout", handleLogout).Use(Jwt())
		auth.POST("/register", handleRegistreUser)
		auth.POST("/refresh", handleLogout).Use(Jwt())
	}
}
