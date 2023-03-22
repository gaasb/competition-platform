package competition_platform

import (
	v1 "github.com/gaasb/competition-platform/internal/api/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	service := v1.TournamentService{}
	server := v1.NewServer(service)
	exit := t.Run()
	server.Start()
	os.Exit(exit)
}

func testTournament(t *testing.T) {
	t.Run("POST:Tournament", func(t *testing.T) {

	})
	t.Run("POST:InvalidDataTournament", func(t *testing.T) {

	})
	t.Run("GET:Tournament", func(t *testing.T) {
		resp, _ := http.NewRequest(http.MethodGet, "api/v1/tournament/", nil)
		assert.Equal(t, http.StatusOK, resp.Response.StatusCode)
	})
	t.Run("GET:InvalidTournamentId", func(t *testing.T) {
		resp, _ := http.NewRequest(http.MethodGet, "api/v1/tournament/", nil)
		assert.Equal(t, http.StatusBadRequest, resp.Response.StatusCode)
	})
	t.Run("GET:TournamentWithBrackets", func(t *testing.T) {
		resp, _ := http.NewRequest(http.MethodGet, "api/v1/tournament?brackets=true", nil)
		assert.Equal(t, http.StatusOK, resp.Response.StatusCode)
	})

	t.Run("PUT:Tournament", func(t *testing.T) {

	})
	t.Run("DELETE:Tournament", func(t *testing.T) {

	})

}
