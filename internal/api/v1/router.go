package v1

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Setup(router *gin.Engine)
}

type TournamentRouter struct{}

func (r *TournamentRouter) Setup(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")

		v1.POST("/login")
		v1.POST("/register")
		v1.GET("/logout")

		tournaments := v1.Group("/tournament/")

		tournaments.OPTIONS("/", handleAllowMethods)      // 405 Method Not Allowed записываем в загаловки и проверяем в другх методах
		tournaments.GET(":id", handleGetTournament)       //✅	?brackets=true
		tournaments.GET("/", handleGetTournamentsList)    //✅	?page prev limit
		tournaments.POST("/", handleNewTournament)        //✅
		tournaments.DELETE(":id", handleDeleteTournament) //✅
		tournaments.PATCH(":id", handleUpdateTournament)  //✅ TODO omit time.Time from JSON and description

		bracket := v1.Group("/bracket")

		bracket.GET(":tid", handleGetBrackets)          //✅
		bracket.POST(":tid", handleNewBracket)          //✅
		bracket.DELETE(":id", handleDeleteBracket)      //✅
		bracket.PATCH(":id", handleUpdateStatusBracket) //✅	?start=true end

		bracket.GET("participants/:id/", handleGetAllParticipants)
		bracket.POST("participants/:id/", handleAddParticipant)
		bracket.DELETE("participants/:id/:team", handleDeleteParticipants)

		match := v1.Group("/match")

		match.POST("/:id")
		match.PUT("/:id", handleUpdateMatch)
		match.DELETE("/:id")

	}
	router.NoRoute(noRoute)

	if gin.Mode() != gin.ReleaseMode {
		// OpenApi Run
	}
}
