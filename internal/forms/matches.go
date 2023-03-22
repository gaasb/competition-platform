package forms

import (
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/volatiletech/null/v8"
)

type Team struct {
	*model.Team
	Participants []*model.Participant
}
type Match struct {
	Round           int        `json:"round" boil:"round"`
	FirstTeam       null.Int64 `json:"first_team" boil:"first_team"`
	SecondTeam      null.Int64 `json:"second_team" boil:"second_team"`
	FirstTeamScore  null.Int64 `json:"first_team_score" boil:"first_team_score"`
	SecondTeamScore null.Int64 `json:"second_team_score" boil:"second_team_score"`
	StartOn         null.Time  `json:"start_on" boil:"start_on"`
	Winner          null.Int64 `json:"winner" boil:"winner"`
}

type MatchForm struct {
	Round           int `form:"round" json:"round"`
	FirstTeamScore  int `form:"first_team_score" json:"first_team_score"`
	SecondTeamScore int `form:"second_team_score" json:"second_team_score"`
}
