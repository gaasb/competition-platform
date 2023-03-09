package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router interface {
	Setup(router *gin.Engine)
}

type TournamentRouter struct {
}

func (r *TournamentRouter) Setup(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.POST("/login")
		v1.POST("/register")
		v1.POST("/logout")

		tournaments := v1.Group("/tournament")
		tournaments.GET("/")
		tournaments.POST("/create")
		tournaments.PUT("/update")
		tournaments.DELETE("/delete")
		tournaments.GET("/:id")
	}
}

type TournamentQuery struct {
	Status string `form:"status"  json:"status" binding:"required"`
	Page   int    `form:"page" json:"page" binding:"required,gte=1,lte=100"`
}

func TournamentsPagination(ctx *gin.Context) {
	ctx.GetInt("start_offset")
	ctx.Param("sort_by")
	//ctx.GetInt()
	var query TournamentQuery
	if err := ctx.ShouldBind(&query); err == nil {
		ctx.JSON(http.StatusOK, gin.H{"result": "Succes", "query": query})
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
}

func TournametsOptions(ctx *gin.Context) {
	session, err := ctx.Cookie("Session")
	ctx.GetHeader("")
	if err != nil {
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "")
		ctx.Abort()
		return
	}
	fmt.Println(session)
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.AbortWithStatus(204)
}
