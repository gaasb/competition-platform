package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ApiServ struct {
	service Service
}

func NewServer(service Service) *ApiServ {
	return &ApiServ{
		service: service,
	}
}

func (s *ApiServ) Start() {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err.Error())
			//os.Exit(1)
		}
	}()

	gratefulShutdown(srv)
}
func gratefulShutdown(srv *http.Server) {
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
