package forms

import (
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
)

var (
	QueryForBrackets = []string{
		model.BracketColumns.ID,
		model.BracketColumns.TypeOf,
		model.BracketColumns.MaxTeams,
		model.BracketColumns.MaxTeamParticipants,
		model.BracketColumns.PlayoffRounds,
		model.BracketColumns.FinalRounds,
		model.BracketColumns.GrandFinalRounds,
		model.BracketColumns.Status,
	}
	selectColumnsForTotalTeams = []string{
		model.BracketTableColumns.ID,
		model.BracketTableColumns.TypeOf,
		model.BracketTableColumns.MaxTeamParticipants,
		model.BracketTableColumns.MaxTeams,
		model.BracketTableColumns.PlayoffRounds,
		model.BracketTableColumns.FinalRounds,
		model.BracketTableColumns.GrandFinalRounds,
		model.BracketTableColumns.Status,
	}
)

type BracketForm struct {
	Type                  model.BracketType `form:"type" json:"type" binding:"required,bracket_type" boil:"type_of"`
	MaxTeam               int               `form:"max_team" json:"max_team" binding:"required,maxteams,lte=128,gt=1" boil:"max_teams"`
	MaxParticipantsInTeam int               `form:"max_participants_in_team" json:"max_participants_in_team" binding:"required,lte=32,gt=0" boil:"max_team_participants"`
	PlayoffRounds         int               `form:"playoff_rounds" json:"playoff_rounds" binding:"required,lte=32,gt=0" boil:"playoff_rounds"`
	FinalRounds           int               `form:"final_rounds" json:"final_rounds" binding:"required,lte=32,gt=0" boil:"final_rounds"`
	GrandFinalRounds      int               `form:"grand_final_rounds" json:"grand_final_rounds" binding:"required,lte=32,gt=0" boil:"grand_final_rounds"`
}

type BracketsWithCount struct {
	model.Bracket `boil:"brackets,bind"`
	TotalTeams    int64 `json:"total_teams" boil:"total_teams"`
}
