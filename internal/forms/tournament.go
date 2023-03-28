package forms

import (
	"fmt"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"

	"github.com/volatiletech/null/v8"
	"time"
)

const (
	BRACKETS_LIMIT    = 5
	BRACKETS_SIZE_MAX = 15
	BRACKETS_SIZE_MIN = 1
)

var (
	QueryForTournament = []string{
		model.TournamentColumns.SportName,
		model.TournamentColumns.Title,
		model.TournamentColumns.Description,
		model.UserAccountColumns.UserLogin,
		model.TournamentColumns.StartAt,
		model.TournamentColumns.EndAt,
		model.TournamentColumns.IsActive,
	}

	InnerJoinForUserName = fmt.Sprintf("%s on %s = %s",
		model.TableNames.UserAccounts,
		model.UserAccountTableColumns.ID,
		model.TournamentTableColumns.CreatedByUser)
)

type TournamentsForm struct {
	Id             string    `form:"-" json:"id" boil:"id"`
	DisciplineName string    `form:"discipline_name" json:"discipline_name" binding:"required,lte=32,gte=2" boil:"sport_name"`
	Title          string    `form:"title" json:"title" binding:"required,lte=32,gte=2" boil:"title"`
	StartAt        time.Time `form:"start_at" json:"start_at" binding:"required_with=EndAt,checktime" boil:"start_at"`
	EndAt          time.Time `form:"end_at" json:"end_at" binding:"required_with=StartAt,gtfield=StartAt" boil:"end_at"`
	Description    null.JSON `form:"description" json:"description" binding:"json" boil:"description"`
}
type TournamentsUpdateForm struct {
	DisciplineName string    `form:"discipline_name" json:"discipline_name,omitempty" binding:"omitempty,lte=32,gte=2"`
	Title          string    `form:"title" json:"title,omitempty" binding:"omitempty,lte=32,gte=2"`
	StartAt        time.Time `form:"start_at" json:"start_at,omitempty" binding:"required_with=EndAt"`
	EndAt          time.Time `form:"end_at" json:"end_at,omitempty" binding:"required_with=StartAt,omitempty,gtfield=StartAt"`
	Description    null.JSON `form:"description" json:"description,omitempty" binding:"json"`
}

type Tournament struct {
	TournamentsForm `boil:"tournaments,bind"`
	CreatedByUser   string           `boil:"user_accounts.user_login" json:"creator"`
	IsActive        bool             `boil:"is_active" json:"is_active"`
	Brackets        []*model.Bracket `boil:"-" json:"brackets,omitempty"`
}
