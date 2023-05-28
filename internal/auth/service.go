package auth

import (
	"database/sql"

	"github.com/gaasb/competition-platform/internal/forms"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

type AuthService struct {
	Service
}

type Service interface {
	RegistreUser(form *forms.RegistrationForm, ctx *gin.Context) error
	AuthUser(form *forms.AuthForm, ctx *gin.Context) error
	UpdateToken(form *User, ctx *gin.Context) error
	DeleteToken(form *User, ctx *gin.Context) error
}
