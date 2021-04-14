package app

import (
	"context"
	agw "github.com/corinz/apigateway/pkg/apigateway"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type app struct {
	router *mux.Router
	server *http.Server
	apis   *agw.APIs
}

func NewApp() *app {
	r := mux.NewRouter().StrictSlash(true)
	s := &http.Server{Handler: r}
	a := &agw.APIs{}
	return &app{router: r, server: s, apis: a}
}

func (a *app) Startup() {
	if _, err := os.Stat("cmd/server/apigateway.json"); err == nil { // save file exists
		if err := a.load(); err != nil {
			log.Println(err.Error())
		}
	}
	a.setupRoutes()
	log.Fatal(a.server.ListenAndServeTLS(".cert/localhost.crt", ".cert/localhost.key"))
}

//TODO review this
func (a *app) Shutdown() {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := a.server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
}
