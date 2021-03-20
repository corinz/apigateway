package main

import (
	"log"
	"net/http"

	agw "github.com/corinz/apigateway/pkg/apigateway"
)

func Startup() {
	apiGW := agw.NewAPIGateway()

	// GETs
	apiGW.R.HandleFunc("/", apiGW.ListAPIs).Methods("GET")
	apiGW.R.HandleFunc("/{api}", apiGW.ListAPI).Methods("GET")
	apiGW.R.HandleFunc("/{api}/{endpoint}", apiGW.ListAPIEndpoints).Methods("GET")

	// POSTs
	apiGW.R.HandleFunc("/", apiGW.CreateAPI).Methods("POST")
	apiGW.R.HandleFunc("/{api}", apiGW.CreateAPIEndpoint).Methods("POST")
	apiGW.R.HandleFunc("/{api}/{endpoint}", apiGW.ExecuteAPIEndpoint).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", apiGW.R))
}

func main() {
	Startup()
}
