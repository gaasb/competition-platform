package v1

import (
	"github.com/gaasb/competition-platform/internal/auth"
	"github.com/gin-gonic/gin"
)

type Router interface {
	Setup(router *gin.Engine)
}

type TournamentRouter struct{}

func (r *TournamentRouter) Setup(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.OPTIONS("/", handleAllowMethods) // 405 Methods Not Allowed если юзер не авторизован

		v1 := api.Group("/v1")
		v1.Use(auth.Jwt())

		tournaments := v1.Group("/tournament/")

		tournaments.GET(":id", handleGetTournament)       //✅	?brackets=true
		tournaments.GET("/", handleGetTournamentsList)    //✅	?page prev limit
		tournaments.POST("/", handleNewTournament)        //✅
		tournaments.DELETE(":id", handleDeleteTournament) //✅
		tournaments.PUT(":id", handleUpdateTournament)    //✅ TODO omit time.Time from JSON and description

		bracket := v1.Group("/bracket")

		bracket.GET(":tid", handleGetBrackets)          //✅
		bracket.POST(":tid", handleNewBracket)          //✅
		bracket.DELETE(":id", handleDeleteBracket)      //✅
		bracket.PATCH(":id", handleUpdateStatusBracket) //✅	?start=true end

		bracket.GET("participants/:id/", handleGetAllParticipants)         //✅
		bracket.POST("participants/:id/", handleAddParticipant)            //✅			?update=true
		bracket.DELETE("participants/:id/:team", handleDeleteParticipants) //✅

		match := v1.Group("/match")

		match.GET(":bid", handleGetAllMatches) //✅
		match.PATCH(":bid", handleUpdateMatchScore)

	}
	router.NoRoute(noRoute)
}
