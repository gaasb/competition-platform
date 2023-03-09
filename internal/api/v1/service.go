package v1

import (
	"context"
	"github.com/gaasb/competition-platform/internal/forms"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
)

type Service interface {
	AddTournament(form forms.TournamentsForm, ctx context.Context) error
	FindAllTournaments(ctx context.Context) ([]*model.Tournament, error)
	GetTournamentBy(id string, ctx context.Context) (*model.Tournament, error)
	UpdateTournamentBy(id string, ctx context.Context) error
	DeleteTournamentBy(id string, ctx context.Context) error

	AddBracket(tournamentId string, form forms.BracketForm, ctx context.Context) error
	FindAllBracketsFrom(tournamentId string, ctx context.Context) []*model.Bracket

	FindAllParticipantFrom(bracketId string, ctx context.Context) []*model.Team
	AddTeamTo(bracketId string, teamAlias string, ctx context.Context) ([]*model.Team, error)
	AddParticipantTo(bracketId string, form forms.ParticipantForm, ctx context.Context) error
	DeleteParticipantBy(bracketId string, form forms.ParticipantForm, ctx context.Context)

	UpdateMatchBy(id string, form forms.MatchForm, ctx context.Context)
}

type TournamentService struct {
}
