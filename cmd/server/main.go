package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	agw "github.com/corinz/apigateway/pkg/apigateway"
)

type app struct {
	serv *http.Server
}

func (a *app) Startup() {
	apiGW := agw.NewAPIGateway()
	a.serv = &http.Server{Addr: ":8080", Handler: apiGW.R}
	defer log.Fatal(a.serv.ListenAndServe())

	// TODO add method to add new routes
	// TODO move handler to the app struct
	// GETs
	apiGW.R.HandleFunc("/", apiGW.ListAPIs).Methods("GET")
	apiGW.R.HandleFunc("/{api}", apiGW.ListAPI).Methods("GET")
	apiGW.R.HandleFunc("/{api}/{endpoint}", apiGW.ListAPIEndpoints).Methods("GET")

	// POSTs
	apiGW.R.HandleFunc("/", apiGW.CreateAPI).Methods("POST")
	apiGW.R.HandleFunc("/{api}", apiGW.CreateAPIEndpoint).Methods("POST")
	apiGW.R.HandleFunc("/{api}/{endpoint}", apiGW.ExecuteAPIEndpoint).Methods("POST")
}

//TODO review this
func (a *app) Shutdown() {
	//a.serv.Shutdown(context.Background())
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := a.serv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
}

func main() {
	app := app{}
	app.Startup()
}
