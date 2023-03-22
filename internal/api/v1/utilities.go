package v1

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
	"github.com/gaasb/competition-platform/internal/utils"
	model "github.com/gaasb/competition-platform/internal/utils/boiler-models"
	"log"
	"time"
)

// Helper function for checking user rights on bracket model
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
		return errors.New("cannot update fields because bracket is finished")
	}
	return nil

}

// Helper function for checking user rights on tournament model without checking status
func verificationForBracketW(userId int64, bracket *model.Bracket, ctx context.Context) error {
	var err error
	var tournament *model.Tournament

	if tournament, err = bracket.Tournament().One(ctx, db); err != nil {
		return errors.New("user verification error")
	}

	if tournament == nil || tournament.CreatedByUser.Int64 != userId {
		return errors.New(utils.NoPermission)
	}

	return nil

}

// Helper function for checking user rights on tournament model
func verificationForTournament(userId int64, tournament *model.Tournament) error {
	if tournament == nil || tournament.CreatedByUser.Int64 != userId {
		return errors.New(utils.NoPermission)
	}

	if tournament.EndAt.Before(time.Now().UTC()) {
		return errors.New("tournament is ended")
	}
	return nil
}

// Helper function for checking user rights
func verificationUserPermission(ctx context.Context) (int64, error) {
	var userId int64
	if value := ctx.Value("user_id"); value != nil {
		if validId, ok := value.(int64); ok {
			userId = validId
		} else {
			return userId, errors.New(utils.NoPermission)
		}
	} else {
		return userId, errors.New(utils.NoPermission)
	}
	return userId, nil
}

func genSingleElimination(totalTeams, toPlayoff, rounds, remain int, bracket *model.Bracket, teams []*model.Team, tx *sql.Tx) error {

	var err error
	r := 0

	for i := 0; i < totalTeams; i++ {

		if (toPlayoff > 0) && (rounds != toPlayoff) {

			if remain < 1 && toPlayoff > 1 {

				if _, err = tx.Exec(
					`INSERT INTO matches (bracket_id, round, first_team, second_team) VALUES ($1, $2, $3, $4)`,
					bracket.ID, rounds+r, teams[i].ID, teams[i+1].ID,
				); err != nil {
					log.Println(err.Error())
					return errors.New("cant add match")
				}
				i += 1
				r += 1
				continue
			}

			if _, err = tx.Exec(
				`INSERT INTO matches (bracket_id, round, first_team) VALUES ($1, $2, $3)`,
				bracket.ID, rounds+r, teams[i].ID,
			); err != nil {
				log.Println(err.Error())
				return errors.New("cant add match")
			}

			if remain > 0 {
				i += 1
			}
			toPlayoff -= 1
		}

		if remain > 0 {

			if _, err = tx.Exec(
				`INSERT INTO matches (bracket_id, round, first_team, second_team) VALUES ($1, $2, $3, $4)`,
				bracket.ID, r, teams[i].ID, teams[i+1].ID,
			); err != nil {
				log.Println(err.Error())
				return errors.New("cant add match")
			}
			i += 1
			r += 1
			remain -= 2
		}
	}
	return err
}
