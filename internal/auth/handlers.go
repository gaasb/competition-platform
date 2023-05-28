package auth

import (
	"net/http"

	"github.com/gaasb/competition-platform/internal/forms"
	"github.com/gaasb/competition-platform/internal/utils"
	"github.com/gin-gonic/gin"
)

// TODO
var service = AuthService{}

const (
	OK  = "ok"
	ERR = "error"
)

// Helper for response
func respMsg(status string, data any, message ...string) gin.H {
	return map[string]any{
		"status":  status,
		"message": message,
		"data":    data,
	}
}

func handleRegistreUser(ctx *gin.Context) {

	var form *forms.RegistrationForm
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		vErr := utils.ValidateErrors(err)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, vErr...))
		return
	}

	if err = service.RegistreUser(form, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, form, "user created"))

}

func handleLogin(ctx *gin.Context) {

	var form *forms.AuthForm
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		vErr := utils.ValidateErrors(err)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, vErr...))
		return
	}

	if err = service.AuthUser(form, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, "successful log in"))

}

func handleLogout(ctx *gin.Context) {

	// var userId int64
	// var err error

	// rawId, hasValue := ctx.Get("user_id")
	// if !hasValue {
	// 	ctx.AbortWithStatusJSON(
	// 		http.StatusBadRequest,
	// 		respMsg(ERR, nil, ""))
	// 	return
	// }
	// if value, ok := rawId.(int64); ok {
	// 	userId = value
	// } else {
	// 	return errors.New(utils.NoPermission)
	// }
	// service.DeleteToken()
	// ctx.JSON(
	// 	http.StatusOK,
	// 	respMsg(OK, nil, "successful log out"))
}

func handleRefreshToken(ctx *gin.Context) {

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, "token has been updated"))
}
