package forms

import (
	"context"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

const (
	BRACKETS_LIMIT = 5
)

type TournamentsForm struct {
	DisciplineName string `form:"sport_name" binding:"required,lte=5" boil:"sport_name"`
	Title          string `form:"required,gte=1,lte=30" json:"title" boil:"title"`
}
type TournamentsWithCounter struct {
	model.Tournament `boil:"bind="`
	Count            int
}

func AddTournament(data model.Tournament, ctx context.Context) {

	data = model.Tournament{
		ID:            uuid.NewV4().String(),
		SportName:     "",
		Title:         "",
		StartAt:       time.Time{},
		EndAt:         time.Time{},
		Description:   null.JSON{},
		CreatedByUser: null.Int64{},
		BracketsLimit: BRACKETS_LIMIT,
		R:             nil,
	}
	uuid.NewV4().String()
	model.Tournament.
		Insert(data, ctx, utils.DB(), boil.Infer())
}
func GetTournament() {}
func (t *TournamentsForm) EditTournament(id string, editable map[string]interface{}) {
	_ = map[string]interface{}{
		model.TournamentColumns.ID: id,
	}
}
