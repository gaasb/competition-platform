package v1

import (
	"fmt"
	"github.com/gaasb/competition-platform/internal/forms"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strconv"
)

var service Service = TournamentService{}

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

// Default error handler
func noRoute(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest,
		respMsg(ERR, nil, "bad request"))
}

func handleAllowMethods(ctx *gin.Context) {
	login := false

	if login {
		ctx.Header("Allow", "GET, POST, PUT, DELETE")

		return
	}
	ctx.Header("Allow", "GET")
}

func handleNewTournament(ctx *gin.Context) {

	var form *forms.TournamentsForm
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		vErr := utils.ValidateErrors(err)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, vErr...))
		return
	}

	//var tournament *forms.TournamentsForm

	if err = form.Create(db, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}
	ctx.JSON(
		http.StatusOK,
		respMsg(OK, form, "created"))
}

func handleDeleteTournament(ctx *gin.Context) {

	idParam := ctx.Param("id")
	tournamentId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}

	err = service.DeleteTournamentBy(tournamentId.String(), ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, fmt.Sprintf("tournament with ID: %s succesfuly deleted", tournamentId)))
}

func handleGetTournamentsList(ctx *gin.Context) {

	currentPage, _ := strconv.Atoi(ctx.Query("page"))
	prevPage, _ := strconv.Atoi(ctx.Query("prev"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	ctx.Set("page", currentPage)
	ctx.Set("prev", prevPage)
	ctx.Set("limit", limit)

	tournaments, err := service.FindAllTournaments(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}
	ctx.JSON(
		http.StatusOK,
		respMsg(OK, tournaments))
}

func handleGetTournament(ctx *gin.Context) {

	withBrackets, err := strconv.ParseBool(ctx.Query("brackets"))
	ctx.Set("brackets", withBrackets)

	idParam := ctx.Param("id")
	tournamentId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}

	tournament, err := service.GetTournamentBy(tournamentId.String(), ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, tournament))
}

func handleUpdateTournament(ctx *gin.Context) {

	idParam := ctx.Param("id")
	tournamentId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}

	var tournament forms.TournamentsUpdateForm
	if err = ctx.ShouldBindJSON(&tournament); err != nil {
		vErr := utils.ValidateErrors(err)

		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, vErr...))
		return
	}

	if err = service.UpdateTournamentBy(tournamentId.String(), tournament, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}
	ctx.JSON(
		http.StatusOK,
		respMsg(OK, tournament, "fields have been updated"))
}

func handleGetBrackets(ctx *gin.Context) {

	idParam := ctx.Param("tid")
	tournamentId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}

	var brackets []*forms.BracketsWithCount
	if brackets, err = service.FindAllBracketsFrom(tournamentId.String(), ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, brackets))
}

func handleDeleteBracket(ctx *gin.Context) {

	idParam := ctx.Param("id")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}
	if err = service.DeleteBracket(bracketId.String(), ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}
	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, fmt.Sprintf("bracket with ID: %s succesfuly deleted", bracketId)))
}

func handleNewBracket(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("tid")
	tournamentId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"))
		return
	}

	var form forms.BracketForm
	if err = ctx.ShouldBindJSON(&form); err != nil {
		vErr := utils.ValidateErrors(err)

		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, vErr...))
		return
	}

	var output *model.Bracket
	if output, err = service.AddBracket(tournamentId.String(), form, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, output))
}

func handleUpdateStatusBracket(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("id")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}

	start, _ := strconv.ParseBool(ctx.Query("start"))
	end, _ := strconv.ParseBool(ctx.Query("end"))

	switch {
	case start:
		err = service.UpdateBracketStatus(bracketId.String(), model.BracketStatusLive, ctx)
	case end:
		err = service.UpdateBracketStatus(bracketId.String(), model.BracketStatusFinished, ctx)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, "status successfully updated"))
}

func handleGetAllParticipants(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("id")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}
	participants, err := service.FindAllParticipantFrom(bracketId.String(), ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()),
		)
		return
	}
	ctx.JSON(
		http.StatusOK,
		respMsg(OK, participants),
	)

}

func handleAddParticipant(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("id")
	isUpdateQuery := ctx.Query("update")
	isUpdate, err := strconv.ParseBool(isUpdateQuery)
	ctx.Set("update", isUpdate)

	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}

	var form forms.ParticipantForm
	if err = ctx.ShouldBindJSON(&form); err != nil {
		vErr := utils.ValidateErrors(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, respMsg(ERR, nil, vErr...))
		return
	}

	if err = service.AddParticipantTo(bracketId.String(), form, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()),
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, form, "participants and team successful added"))
}

func handleDeleteParticipants(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("id")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}
	var teamParam string
	if teamParam = ctx.Param("team"); len(teamParam) < 2 {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "incorrect team parameter value"),
		)
		return
	}

	if err = service.DeleteParticipantsBy(bracketId.String(), teamParam, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, "team with participants successfully deleted"))
}

func handleUpdateMatchScore(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("bid")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}

	var form forms.MatchForm
	if err = ctx.ShouldBindJSON(&form); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, utils.ValidateErrors(err)...))
		return
	}

	if err = service.UpdateMatchScoreBy(bracketId.String(), form, ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()))
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, nil, "match successfully updated"))
}

func handleGetAllMatches(ctx *gin.Context) {

	var err error
	idParam := ctx.Param("bid")
	bracketId, err := uuid.FromString(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, "id must be in uuid format"),
		)
		return
	}

	var matches []*forms.Match
	if matches, err = service.FindAllMatches(bracketId.String(), ctx); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			respMsg(ERR, nil, err.Error()),
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		respMsg(OK, matches))
}
