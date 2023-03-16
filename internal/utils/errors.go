package utils

import (
	"log"
)

const (
	NoPermission  = "You dont have permission to this operation."
	Unauthorized  = "Unauthorized"
	NoTournament  = "There are no tournaments."
	NoBrackets    = "There are no brackets."
	NoMatches     = "There are no matches."
	BeginTxErr    = "Unable to begin transaction."
	RollbackTxErr = "Unable to rollback transaction."
	CommitTxErr   = "Unable to commit transaction."
)

func dieIf(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
