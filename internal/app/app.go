package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	agw "github.com/corinz/apigateway/pkg/apigateway"
	"github.com/gorilla/mux"
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

func (a *app) Startup(addr string) {
	a.setupRoutes()
	a.server.Addr = addr
	log.Fatal(a.server.ListenAndServe())
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