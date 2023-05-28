package v1

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gaasb/competition-platform/internal/auth"
	"github.com/gaasb/competition-platform/internal/utils"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	service Service
	router  Router
}

func NewServer(service Service, router Router) *ApiServer {

	return &ApiServer{
		service: service,
		router:  router,
		// router:  &TournamentRouter{},
	}

}

func (s *ApiServer) Start() {

	utils.Init()
	auth.InitJwtSecret()
	utils.SetupValidator()
	db = utils.GetDB()
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	s.router.Setup(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	gratefulShutdown(srv)
}

func gratefulShutdown(srv *http.Server) {

	defer utils.CloseDB()

	closeServer := make(chan os.Signal, 1)
	signal.Notify(closeServer, os.Interrupt)

	<-closeServer
	log.Println(">>Closing server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	<-ctx.Done()
	log.Println(">>Server closed")
}
