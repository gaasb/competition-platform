package main

import (
	v1 "github.com/gaasb/competition-platform/internal/api/v1"
)

func main() {

	service := v1.TournamentService{}
	server := v1.NewServer(service, &v1.TournamentRouter{})

	server.Start()
}
