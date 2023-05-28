package main

import (
	v1 "github.com/gaasb/competition-platform/internal/api/v1"
	"github.com/gaasb/competition-platform/internal/auth"
)

func main() {

	service := auth.AuthService{}
	server := v1.NewServer(service, &auth.AuthRouter{})

	server.Start()
}
