package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gaasb/competition-platform/internal/forms"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"math"
	"math/rand"
	"time"
)

var db *sql.DB

type TournamentService struct{}

type Service interface {
	AddTournament(form *forms.TournamentsForm, ctx *gin.Context) error
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

	FindAllMatches(bracketId string, ctx context.Context) ([]*forms.Match, error)
	GenEliMatch(bracket *model.Bracket, ctx context.Context) error
	UpdateMatchScoreBy(bracketId string, form forms.MatchForm, ctx context.Context) error
}

func (t TournamentService) AddTournament(form *forms.TournamentsForm, ctx *gin.Context) error {

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
	form.Id = id.String()

	tournament := model.Tournament{
		ID:            form.Id,
		SportName:     form.DisciplineName,
		Title:         form.Title,
		StartAt:       form.StartAt.UTC(),
		EndAt:         form.EndAt.UTC(),
		Description:   form.Description,
		CreatedByUser: null.Int64From(userId),
		BracketsLimit: forms.BRACKETS_LIMIT,
	}

	if err = tournament.Insert(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return errors.New("invalid data")
	}

	return nil

}

func (t TournamentService) FindAllTournaments(ctx context.Context) (any, error) {

	currentPage := ctx.Value("page").(int)
	prevPage := ctx.Value("prev").(int)
	limit := ctx.Value("limit").(int)

	switch {
	case limit == 0:
		limit = forms.BRACKETS_SIZE_MIN
	case limit < forms.BRACKETS_SIZE_MIN:
		return nil, errors.New(fmt.Sprintf("limit cannot be less than 1 %d", forms.BRACKETS_SIZE_MAX))
	case limit > forms.BRACKETS_SIZE_MAX:
		return nil, errors.New(fmt.Sprintf("limit cannot be greater than %d", forms.BRACKETS_SIZE_MAX))
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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var tournament *model.Tournament
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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var tournament *model.Tournament
	tournament, err = model.FindTournament(ctx, db, id)
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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return nil, err
	}

	var tournament *model.Tournament
	tournament, err = model.FindTournament(ctx, db, tournamentId)
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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var bracket *model.Bracket

	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}

	switch status {
	case model.BracketStatusFinished:
		if bracket.Status == model.BracketStatusFinished {
			return errors.New(fmt.Sprintf("current status is %s", bracket.Status.String()))
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
			if err = t.GenEliMatch(bracket, ctx); err != nil {
				return err
			}
		} else {
			return errors.New("the current number of teams is less than 2")
		}
	default:
		return errors.New("unknown status")
	}

	if _, err = bracket.Update(ctx, db, boil.Infer()); err != nil {
		log.Println(err.Error())
		return errors.New("failed on update status")
	}
	return nil

}

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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracketW(userId, bracket, ctx); err != nil {
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

	var teams []*model.Team
	if teams, err = bracket.Teams().All(ctx, db); err != nil {
		return nil, errors.New("teams is empty")
	}

	output := make([]forms.ParticipantsFromTeam, len(teams))

	for i, item := range teams {

		var participants []*forms.Participant
		if err = item.Participants().Bind(ctx, db, &participants); err != nil {
			log.Println(err.Error())
			return nil, errors.New("participants not found")
		}
		output[i].Team = item.TeamAlias
		output[i].Participants = participants

	}

	if err != nil {
		log.Println(err)
		return nil, errors.New("teams not found")
	}
	return output, nil
}

func (t TournamentService) AddTeamTo(bracketId string, teamAlias string, ctx context.Context) (*model.Team, error) {

	var userId int64
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return nil, err
	}

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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var isUpdate bool
	if value, ok := ctx.Value("update").(bool); ok {
		isUpdate = value
	}

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

	if isUpdate {

		var team *model.Team
		if team, err = bracket.Teams(model.TeamWhere.TeamAlias.EQ(form.Team)).One(ctx, db); err != nil {
			return errors.New("team not found")
		}

		var participants model.ParticipantSlice
		if participants, err = team.Participants(qm.Limit(bracket.MaxTeamParticipants)).All(ctx, db); err != nil {
			return errors.New("participants is empty")
		}

		if len(participants) != len(form.Participants) {
			return errors.New("the number of participants should be equal to bracket max participants")
		}

		for i, _ := range participants {

			participants[i].UserAlias = form.Participants[i].UserAlias
			participants[i].Contact = null.StringFrom(form.Participants[i].Contact)

			if _, err = participants[i].Update(ctx, db, boil.Infer()); err != nil {
				return errors.New("problem to update participant")
			}
		}

		return nil
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
		form.Team, bracketId).
		Scan(&teamId); err != nil || isUpdate {
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
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId); err != nil {
		return errors.New("no brackets yet")
	}

	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return errors.New(utils.NoPermission)
	}

	if bracket.Status == model.BracketStatusLive {
		return errors.New("cant delete because bracket is live")
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

func (t TournamentService) FindAllMatches(bracketId string, ctx context.Context) ([]*forms.Match, error) {

	var err error

	var bracket *model.Bracket
	if bracket, err = model.FindBracket(ctx, db, bracketId, model.BracketColumns.ID); err != nil {
		return nil, errors.New("bracket not found")
	}

	var matches []*forms.Match
	if err = queries.Raw(`select m.round, t1.team_alias first_team, t2.team_alias second_team, m.first_team_score, m.second_team_score, m.start_on, t3.team_alias team_winner
from matches m 
    LEFT OUTER JOIN teams t1 on t1.id = m.first_team 
    LEFT OUTER JOIN teams t2 on t2.id = m.second_team
    LEFT OUTER JOIN teams t3 on t3.id = m.winner
where m.bracket_id = $1;`, bracket.ID).
		Bind(ctx, db, &matches); err != nil {
		log.Println(err.Error())
		return nil, errors.New("matches not found")
	}

	return matches, nil
}

func (t TournamentService) UpdateMatchScoreBy(bracketId string, form forms.MatchForm, ctx context.Context) error {

	var userId int64
	var err error

	if userId, err = verificationUserPermission(ctx); err != nil {
		return err
	}

	var match *model.Match
	if match, err = model.Matches(model.MatchWhere.BracketID.EQ(null.StringFrom(bracketId)), model.MatchWhere.Round.EQ(form.Round)).One(ctx, db); err != nil {
		return errors.New("no match yet")
	}

	if form.FirstTeamScore == form.SecondTeamScore {
		return errors.New("the score should not to be the same")
	}

	if match.FirstTeam.IsZero() || match.SecondTeam.IsZero() {
		return errors.New("one of the team field is empty")
	}

	var bracket *model.Bracket
	if bracket, err = match.Bracket().One(ctx, db); err != nil {
		return errors.New(utils.NoPermission)
	}
	if err = verificationForBracket(userId, bracket, ctx); err != nil {
		return err
	}

	var totalTeams int64
	if totalTeams, err = bracket.Teams().Count(ctx, db); err != nil {
		return errors.New("some problem to count total team")
	}

	rounds := int(math.Pow(2, math.Floor(math.Log2(float64(totalTeams)))))

	match.FirstTeamScore = null.IntFrom(form.FirstTeamScore)
	match.SecondTeamScore = null.IntFrom(form.SecondTeamScore)

	var loserTeam int64
	var winnerTeam int64
	//set winner
	if form.FirstTeamScore > form.SecondTeamScore {
		winnerTeam = match.FirstTeam.Int64
		loserTeam = match.SecondTeam.Int64
	} else {
		winnerTeam = match.SecondTeam.Int64
		loserTeam = match.FirstTeam.Int64
	}
	match.Winner = null.Int64From(winnerTeam)

	if _, err = match.Update(ctx, db, boil.Infer()); err != nil {
		return errors.New("failed on update match")
	}

	currentRound := int(math.Abs(float64(match.Round)))

	if currentRound >= rounds*2-2 {
		return nil
	}

	var nextRound int
	var position int
	position = currentRound - rounds
	position -= position % 2 //if position%2 != 0 { position -= 1 }
	position = position / 2

	a := rounds - (int(totalTeams) - rounds)
	b := int(totalTeams) - rounds
	if currentRound < int(math.Min(float64(a), float64(b))) {
		nextRound = currentRound + rounds
	} else {
		nextRound = ((rounds / 2) + position) + rounds
	}

	if match.Round < 0 {
		nextRound = 0 - nextRound
	}

	pos := match.Round

	var nextMatch *model.Match
	if nextMatch, err = bracket.Matches(model.MatchWhere.Round.EQ(nextRound)).One(ctx, db); err != nil {

		nextMatch = &model.Match{
			BracketID: null.StringFrom(bracket.ID),
			Round:     nextRound,
		}

		if match.Round < rounds {
			pos += 1
		}

		if pos%2 == 0 {
			nextMatch.FirstTeam = match.Winner
		} else {
			nextMatch.SecondTeam = match.Winner
		}

		if err = nextMatch.Insert(ctx, db, boil.Infer()); err != nil {
			return errors.New("some problem to link next match")
		}

		if bracket.TypeOf != model.BracketTypeDOUBLE_ELIMINATION {
			return nil
		}

	}
	if match.Round < rounds {
		pos += 1
	}
	if pos%2 == 0 {
		nextMatch.FirstTeam = match.Winner
	} else {
		nextMatch.SecondTeam = match.Winner
	}

	if !nextMatch.Winner.IsZero() {
		return errors.New("cant update because in linked match winner is already")
	}

	if _, err = nextMatch.Update(ctx, db, boil.Infer()); err != nil {
		return errors.New("some problem to update next match")
	}

	if bracket.TypeOf == model.BracketTypeDOUBLE_ELIMINATION && match.Round >= 0 {

		var nextLoserMatch *model.Match
		if nextLoserMatch, err = bracket.Matches(model.MatchWhere.Round.EQ(nextRound)).One(ctx, db); err != nil {

			nextLoserMatch = &model.Match{
				BracketID: null.StringFrom(bracket.ID),
				Round:     nextRound,
			}

			if position%2 == 0 {
				nextLoserMatch.FirstTeam = null.Int64From(loserTeam)
			} else {
				nextLoserMatch.SecondTeam = null.Int64From(loserTeam)
			}

			if err = nextLoserMatch.Insert(ctx, db, boil.Infer()); err != nil {
				return errors.New("some problem to link next loser match")
			}
			return nil
		}

		if !nextLoserMatch.Winner.IsZero() {
			return errors.New("cant update because in linked match winner is already")
		}

		if position%2 == 0 {
			nextLoserMatch.FirstTeam = null.Int64From(loserTeam)
		} else {
			nextLoserMatch.SecondTeam = null.Int64From(loserTeam)
		}

		if _, err = nextLoserMatch.Update(ctx, db, boil.Infer()); err != nil {
			return errors.New("some problem to update next match")
		}
	}

	if nextRound == rounds*2-2 {
		//TODO конец сетки последний лузер вс последний винер
	}

	return nil
}

func (t TournamentService) GenEliMatch(bracket *model.Bracket, ctx context.Context) error {

	var err error
	//No verification because it is meant to be used in the method UpdateBracketStatus on start query.
	var teams []*model.Team
	if teams, err = bracket.Teams().All(ctx, db); err != nil {
		log.Println(err.Error())
		return errors.New("teams not found")
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(teams), func(i, j int) {
		teams[i], teams[j] = teams[j], teams[i]
	})

	var rounds int
	totalTeams := len(teams)
	rounds = int(math.Pow(2, math.Floor(math.Log2(float64(totalTeams)))))
	remain := (totalTeams - rounds) * 2
	toPlayoff := totalTeams - remain

	if remain == 0 {
		remain = rounds
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		log.Println(err.Error())
		return errors.New("some problems")
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Println(err.Error())
				err = errors.New("some problems")
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Println(err.Error())
			err = errors.New("some problems")
		}
	}()

	switch bracket.TypeOf {
	case model.BracketTypeSINGLE_ELIMINATION:

		if err = genSingleElimination(totalTeams, toPlayoff, rounds, remain, bracket, teams, tx); err != nil {
			return err
		}

	case model.BracketTypeDOUBLE_ELIMINATION:

		if err = genSingleElimination(totalTeams, toPlayoff, rounds, remain, bracket, teams, tx); err != nil {
			return err
		}

	case model.BracketTypeROUND_ROBIN:

	default:
		return errors.New("invalid bracket type")
	}

	return err
}
