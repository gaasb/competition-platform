package forms

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/gin-gonic/gin"
	"log"

	//"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (t *TournamentsForm) Create(db *sql.DB, ctx *gin.Context) error {

	var userId int64
	rawId, hasValue := ctx.Get("user_id")
	if !hasValue {
		return errors.New(utils.NoPermission)
	}
	if value, ok := rawId.(int64); ok {
		userId = value
	} else {
		return errors.New(utils.NoPermission)
	}

	id, err := uuid.NewV6()
	if err != nil {
		log.Println(err.Error())
		return errors.New("problem on generation uuid")
	}
	t.Id = id.String()

	tournament := model.Tournament{
		ID:            t.Id,
		SportName:     t.DisciplineName,
		Title:         t.Title,
		StartAt:       t.StartAt.UTC(),
		EndAt:         t.EndAt.UTC(),
		Description:   t.Description,
		CreatedByUser: null.Int64From(userId),
		BracketsLimit: BRACKETS_LIMIT,
	}

	if err = tournament.Insert(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return errors.New("invalid data")
	}

	return nil

}
