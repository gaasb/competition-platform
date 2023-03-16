package forms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"math"
	"time"
)

type Team struct {
	*model.Team
	Participants []*model.Participant
}
type Match struct {
	Round           int
	FirstTeam       Team
	SecondTeam      Team
	FirstTeamScore  int
	SecondTeamScore int
	StartOn         time.Time
	Winner          *Team
}

type MatchForm struct {
	FirstTeamScore  int
	SecondTeamScore int
}

func (m *MatchForm) Create(ctx context.Context, db *sql.DB, bracketId string, teams []*model.Team) error {
	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}
	// Verification
	var rounds int
	totalTeams := len(teams)
	rounds = int(math.Pow(2, math.Round(math.Log2(float64(totalTeams)))))
	remaind := rounds - totalTeams

	fmt.Println(remaind)

	switch bracket.TypeOf {
	case model.BracketTypeSINGLE_ELIMINATION:

	case model.BracketTypeDOUBLE_ELIMINATION:

	case model.BracketTypeROUND_ROBIN:

	default:
		return errors.New("invalid bracket type")
	}

	return nil
}
