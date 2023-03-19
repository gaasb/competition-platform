package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gaasb/competition-platform/internal/forms"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"math"
	"time"
)

var db *sql.DB

type TournamentService struct{}

type Service interface {
	FindAllTournaments(ctx context.Context) (any, error)
	GetTournamentBy(id string, ctx context.Context) (*forms.Tournament, error)
	UpdateTournamentBy(id string, params forms.TournamentsUpdateForm, ctx context.Context) error
	DeleteTournamentBy(id string, ctx context.Context) error

	AddBracket(tournamentId string, form forms.BracketForm, ctx context.Context) (*model.Bracket, error)
	DeleteBracket(bracketId string, ctx context.Context) error
	UpdateBracketStatus(bracketId string, status model.BracketStatus, ctx context.Context) error
	FindAllBracketsFrom(tournamentId string, ctx context.Context) ([]*forms.BracketsWithCount, error)

	FindAllParticipantFrom(bracketId string, ctx context.Context) ([]forms.ParticipantsFromTeam, error)
	AddTeamTo(bracketId string, teamAlias string, ctx context.Context) (*model.Team, error)
	AddParticipantTo(bracketId string, form forms.ParticipantForm, ctx context.Context) error
	DeleteParticipantsBy(bracketId string, teamAlias string, ctx context.Context) error

	//AddMatch(bracketId string, form forms.MatchForm, ) error
	UpdateMatchBy(id int64, form forms.MatchForm, ctx context.Context) error
}

func (t TournamentService) FindAllTournaments(ctx context.Context) (any, error) {

	currentPage := ctx.Value("page").(int)
	prevPage := ctx.Value("prev").(int)
	limit := ctx.Value("limit").(int)

	switch {
	case limit == 0:
		limit = forms.BRACKETS_SIZE_MIN
	case limit < forms.BRACKETS_SIZE_MIN:
		return nil, errors.New(fmt.Sprintf("limit cannot be less than 1 %s", forms.BRACKETS_SIZE_MAX))
	case limit > forms.BRACKETS_SIZE_MAX:
		return nil, errors.New(fmt.Sprintf("limit cannot be greater than %s", forms.BRACKETS_SIZE_MAX))
	case currentPage < 0:
		return nil, errors.New("page must be positive")
	case prevPage < 0:
		return nil, errors.New("prev must be positive")
	}

	count, err := model.Tournaments().Count(ctx, db)
	totalCount := int(count)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("some problem to find tournaments")
	}

	orderBy := model.TournamentTableColumns.ID
	if currentPage < prevPage {
		orderBy += ` asc`
	} else {
		orderBy += ` desc`
	}

	offset := currentPage * limit
	maxPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	var output []*forms.Tournament

	if err = model.NewQuery(
		qm.Select(append(forms.QueryForTournament, model.TournamentTableColumns.ID)...),
		qm.From(model.TableNames.Tournaments),
		qm.InnerJoin(forms.InnerJoinForUserName),
		qm.OrderBy(orderBy),
		qm.Offset(offset),
		qm.Limit(limit)).
		Bind(ctx, db, &output); err != nil {

		log.Println(err.Error())
		return nil, errors.New("tournaments not found")
	}

	return struct {
		TotalPages  int                 `json:"total_pages"`
		Tournaments []*forms.Tournament `json:"tournaments"`
	}{maxPages, output}, nil

}

func (t TournamentService) GetTournamentBy(id string, ctx context.Context) (*forms.Tournament, error) {

	var err error
	var output forms.Tournament

	var withBrackets bool
	if value := ctx.Value("brackets"); value != nil {
		withBrackets = value.(bool)
	}

	if err = model.NewQuery(
		qm.Select(forms.QueryForTournament...),
		qm.From(model.TableNames.Tournaments),
		qm.InnerJoin(forms.InnerJoinForUserName),
		model.TournamentWhere.ID.EQ(id),
		qm.Limit(1)).
		Bind(ctx, db, &output); err != nil {

		log.Println(err.Error())
		return nil, errors.New("tournament not found")
	}

	if withBrackets {
		var brackets []*model.Bracket

		if brackets, err = model.Brackets(
			qm.Select(forms.QueryForBrackets...),
			qm.Load(model.BracketRels.Tournament),
			model.BracketWhere.TournamentID.EQ(null.StringFrom(id))).
			All(ctx, db); err != nil {

			log.Println(err.Error())
			return nil, errors.New("brackets not found")
		}
		output.Brackets = brackets
	}
	return &output, nil
}

func (t TournamentService) UpdateTournamentBy(id string, params forms.TournamentsUpdateForm, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var tournament *model.Tournament
	var err error

	if tournament, err = model.FindTournament(ctx, db, id); err != nil {
		return errors.New("tournament not found")
	}

	if tournament.CreatedByUser.Int64 != userId {
		return errors.New(utils.NoPermission)
	}

	timeNow := time.Now().UTC()
	canUpdate := false

	if len(params.Title) > 0 {
		tournament.Title = params.Title
		canUpdate = true
	}
	if len(params.DisciplineName) > 0 {
		tournament.SportName = params.DisciplineName
		canUpdate = true
	}
	if !params.Description.IsZero() && params.Description.Valid && len(params.Description.JSON) > 2 {
		tournament.Description = params.Description
		canUpdate = true
	}
	if params.StartAt.UTC().Before(params.EndAt.UTC()) && (params.StartAt.UTC().After(timeNow)) {
		tournament.StartAt = params.StartAt.UTC()
		canUpdate = true
	}
	if params.EndAt.UTC().After(params.StartAt.UTC()) && (params.EndAt.UTC().After(timeNow)) {
		tournament.EndAt = params.EndAt.UTC()
		canUpdate = true
	}

	if !canUpdate {
		return errors.New("no data for update")
	}

	if _, err = tournament.Update(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return errors.New("failed to update")
	}
	return nil
}

func (t TournamentService) DeleteTournamentBy(id string, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	tournament, err := model.FindTournament(ctx, db, id)
	if err != nil {
		return errors.New("There are no tournaments to delete")
	}
	if tournament.CreatedByUser.Int64 != userId {
		return errors.New(utils.NoPermission)
	}

	_, err = tournament.Delete(ctx, db)
	if err != nil {
		log.Println(err.Error())
		return errors.New("cant on delete tournament")
	}
	return nil
}

func (t TournamentService) AddBracket(tournamentId string, form forms.BracketForm, ctx context.Context) (*model.Bracket, error) {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return nil, errors.New(utils.NoPermission)
		}
	} else {

		return nil, errors.New(utils.NoPermission)
	}

	tournament, err := model.FindTournament(ctx, db, tournamentId)
	if err != nil {
		return nil, errors.New("There are no tournaments")
	}

	if err = verificationForTournament(userId, tournament); err != nil {
		return nil, err
	}

	count, err := tournament.Brackets().Count(ctx, db)
	if err != nil || int64(tournament.BracketsLimit) <= count {
		return nil, errors.New("tournament reach max teams ")
	}

	id, err := uuid.NewV6()
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("problem with generation id for bracket")
	}

	bracket := model.Bracket{
		ID:                  id.String(),
		TypeOf:              form.Type,
		MaxTeams:            null.IntFrom(form.MaxTeam),
		MaxTeamParticipants: form.MaxParticipantsInTeam,
		TournamentID:        null.StringFrom(tournamentId),
		PlayoffRounds:       form.PlayoffRounds,
		FinalRounds:         form.FinalRounds,
		GrandFinalRounds:    form.GrandFinalRounds,
	}

	if err = bracket.Insert(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return nil, errors.New("bracket need to be unique")

	}
	return &bracket, nil
}

func (t TournamentService) UpdateBracketStatus(bracketId string, status model.BracketStatus, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var bracket *model.Bracket
	var err error

	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}

	switch status {
	case model.BracketStatusFinished:
		if bracket.Status == model.BracketStatusFinished {

		}
		bracket.Status = model.BracketStatusFinished

	case model.BracketStatusLive:
		if bracket.Status != model.BracketStatusPending {
			return errors.New(fmt.Sprintf("current status is %s", bracket.Status.String()))
		}

		var totalTeams int64
		totalTeams, err = model.Teams(model.TeamWhere.BracketID.EQ(null.StringFrom(bracketId))).Count(ctx, db)

		if int(totalTeams) >= 2 {
			bracket.Status = model.BracketStatusLive
		} else {
			return errors.New("the current number of teams is less than 2")
		}
	default:
		return errors.New("unknown status")
	}

	if status == model.BracketStatusLive { //TODO CREATE MATCHES
		var totalTeams int64
		totalTeams, err = model.Teams(model.TeamWhere.BracketID.EQ(null.StringFrom(bracketId))).Count(ctx, db)

		if bracket.MaxTeams.Int <= int(totalTeams) {
			bracket.Status = model.BracketStatusLive
		} else {
			return errors.New("the current number of teams is less than expected")
		}
	}

	if _, err = bracket.Update(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return errors.New("failed on update status")
	}
	return nil

} //TODO <---

func (t TournamentService) FindAllBracketsFrom(tournamentId string, ctx context.Context) ([]*forms.BracketsWithCount, error) {

	var err error
	if _, err = model.FindTournament(ctx, db, tournamentId); err != nil {
		return nil, errors.New("tournament not found")
	}

	var brackets []*forms.BracketsWithCount
	err = queries.Raw(
		`select b.id, b.type_of, b.max_teams, b.max_team_participants, b.playoff_rounds, b.final_rounds, b.grand_final_rounds, b.status, count(t.id) as "total_teams" from brackets b left join teams t on t.bracket_id = b.id where b.tournament_id = $1 group by b.id`,
		tournamentId).
		Bind(ctx, db, &brackets)
	if err != nil {
		log.Println(err)
		return nil, errors.New("some problems to find brackets")
	}
	return brackets, nil

}

func (t TournamentService) DeleteBracket(bracketId string, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}

	_, err = bracket.Delete(ctx, db)
	if err != nil {
		return errors.New("Unable to delete bracket.")
	}
	return nil

}

func (t TournamentService) FindAllParticipantFrom(bracketId string, ctx context.Context) ([]forms.ParticipantsFromTeam, error) {

	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return nil, errors.New("no brackets yet")
	}
	//fmt.Println(*bracket.Teams(qm.Select("team_alias"), qm.Load(model.TeamRels.Participants)).)
	var teams []*model.Team
	teams, err = bracket.Teams(qm.Load(model.TeamRels.Participants)).All(ctx, db)

	output := make([]forms.ParticipantsFromTeam, len(teams))
	for i, item := range teams {
		output[i].Team = item.TeamAlias
		output[i].Participants = item.R.Participants
	}

	if err != nil {
		log.Println(err)
		return nil, errors.New("teams not found")
	}
	return output, nil
}

func (t TournamentService) AddTeamTo(bracketId string, teamAlias string, ctx context.Context) (*model.Team, error) {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return nil, errors.New(utils.NoPermission)
		}
	} else {
		return nil, errors.New(utils.NoPermission)
	}

	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return nil, errors.New("no brackets yet")
	}
	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return nil, err
	}

	team := model.Team{
		TeamAlias: teamAlias,
		BracketID: null.StringFrom(bracketId),
	}
	err = team.Insert(ctx, db, boil.Infer())
	return &team, err

}

func (t TournamentService) AddParticipantTo(bracketId string, form forms.ParticipantForm, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var isUpdate bool // TODO <----------------------------------------------------------
	if value, ok := ctx.Value("update").(bool); ok {
		isUpdate = value
	}

	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}
	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}
	if len(form.Participants) != bracket.MaxTeamParticipants {
		return errors.New("the number of participants should be equal to bracket max participants")
	}

	if bracket.Status != model.BracketStatusPending {
		return errors.New(fmt.Sprintf("bracket is %s", bracket.Status.String()))
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("%s %v", utils.BeginTxErr, err)
		return errors.New(utils.BeginTxErr)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Println(err.Error())
			}
		} else {
			err = tx.Commit()
		}
	}()

	var teamId int64
	if err = tx.QueryRow(
		`INSERT INTO teams (team_alias, bracket_id) VALUES ($1, $2) RETURNING id`,
		form.Team,
		bracketId).Scan(&teamId); err != nil || isUpdate {
		return errors.New("team is already in bracket")
	}

	var totalTeams int64
	if totalTeams, err = bracket.Teams().Count(ctx, db); err != nil {
		return errors.New("the number of participants reached a maximum")
	}

	if int(totalTeams) > bracket.MaxTeams.Int {
		err = errors.New("the number of teams reached a maximum")
		return err
	}

	for _, item := range form.Participants {

		if _, err = tx.Exec(
			`INSERT INTO participants (user_alias, team_id, contact) VALUES ($1, $2, $3)`,
			item.UserAlias,
			teamId,
			item.Contact); err != nil {
			return errors.New("participant user alias need to be unique in team")
		}

	}
	return nil
}

func (t TournamentService) DeleteParticipantsBy(bracketId string, teamAlias string, ctx context.Context) error {

	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var err error
	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return errors.New(utils.NoPermission)
	}

	var team *model.Team
	if team, err = bracket.Teams(model.TeamWhere.TeamAlias.EQ(teamAlias)).One(ctx, db); err != nil {
		return errors.New("team not found")
	}

	if _, err = team.Delete(ctx, db); err != nil {
		return errors.New("cant delete team")
	}

	return nil
}

func (t TournamentService) UpdateMatchBy(id int64, form forms.MatchForm, ctx context.Context) error {

	var err error
	var userId int64

	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return errors.New(utils.NoPermission)
		}
	} else {
		return errors.New(utils.NoPermission)
	}

	var match *model.Match
	if match, err = model.FindMatch(ctx, db, id); err != nil {
		return errors.New("no match yet")
	}

	var bracket *model.Bracket
	if bracket, err = match.Bracket().One(ctx, db); err != nil {
		return errors.New(utils.NoPermission)
	}
	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}

	match.FirstTeamScore = null.IntFrom(form.FirstTeamScore)
	match.SecondTeamScore = null.IntFrom(form.SecondTeamScore)

	if form.FirstTeamScore > form.SecondTeamScore {
		match.Winner = match.FirstTeam
	} else {
		match.Winner = match.SecondTeam
	}

	if _, err = match.Update(ctx, db, boil.Infer()); err != nil {
		return errors.New("failed on update match")
	}
	//TODO Link to match with next round or create new match
	return nil
}

func verificationForBracket(userId int64, bracket *model.Bracket, ctx context.Context) error {
	var err error
	var tournament *model.Tournament

	if tournament, err = bracket.Tournament().One(ctx, db); err != nil {
		return errors.New("user verification error")
	}
	if err = verificationForTournament(userId, tournament); err != nil {
		return err
	}

	if bracket.Status.String() == model.BracketStatusFinished.String() {
		return errors.New("cannot change status because bracket is finished")
	}
	return nil

}

func verificationForTournament(userId int64, tournament *model.Tournament) error {
	if tournament == nil || tournament.CreatedByUser.Int64 != userId {
		return errors.New(utils.NoPermission)
	}

	if tournament.EndAt.Before(time.Now().UTC()) {
		return errors.New("tournament is ended")
	}
	return nil
}
