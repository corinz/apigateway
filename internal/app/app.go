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
	s := &http.Server{Addr: ":8080", Handler: r}
	a := &agw.APIs{}
	return &app{router: r, server: s, apis: a}
}

func (a *app) Start(addr string) error {
	a.server.Addr = addr
	return a.server.ListenAndServe()
}

// TODO error check empty or existing name
func (a *app) addSubRouter(api *agw.API) {
	api.Router = a.router.PathPrefix("/" + api.Name).Subrouter() // "/{apiName}/"
}

// TODO error check empty or existing name
func (a *app) newHandleFunc(api *agw.API, subPath string) {
	api.Router.HandleFunc("/"+subPath, Generic) // "/{apiName}/{aepName}/"
}

func (a *app) Startup() {

	// TODO add method to add new routes
	// GETs
	a.router.HandleFunc("/", a.ListAPIs).Methods("GET")
	a.router.HandleFunc("/{api}", a.ListAPI).Methods("GET")
	a.router.HandleFunc("/{api}/{endpoint}", a.ListAPIEndpoints).Methods("GET")

	// POSTs
	a.router.HandleFunc("/", a.CreateAPI).Methods("POST")
	a.router.HandleFunc("/{api}", a.CreateAPIEndpoint).Methods("POST")
	a.router.HandleFunc("/{api}/{endpoint}", a.ExecuteAPIEndpoint).Methods("POST")

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