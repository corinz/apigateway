package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

// apiEndpoint is a struct representing a single API Endpoint with a route and http verb
type apiEndpoint struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	HTTPVerb    string `json:"HttpVerb`
	Route       string `json:"Route`
	Command     string `json:"Command"`
	UID         int
}

// API is a struct representing APIEndpoints
type api struct {
	Name   string `json:"Name"`
	Route  string `json:"Route"`
	apiEPs []apiEndpoint
	router *mux.Router
}

// struct that wraps globals
// instantiate it in main()

// APIs is a global var representing multiple API structs
var apis []api

//
var mainRouter = mux.NewRouter().StrictSlash(true) // make global with intent of attaching new routers

func createApiEP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	apiEP := unmarshalApiEP(r)
	api := getAPI(apiName)

	appendApiEP(apiName, apiEP)
	newHandleFunc(api, apiEP.Name)
	json.NewEncoder(w).Encode(apiEP)
}

func appendApiEP(apiName string, apiEP apiEndpoint) {
	// TODO validate this method
	api := getAPI(apiName)
	api.apiEPs = append(api.apiEPs, apiEP)
}

func createApi(w http.ResponseWriter, r *http.Request) {
	api := unmarshalApi(r) //TODO Test to see if API already exists and validate the Route is unique/valid

	// Initialize api with root endpoint and append
	apiRootEP := apiEndpoint{
		Name:        api.Name + "/rootEP",
		Description: api.Name + "/rootEP",
		HTTPVerb:    "GET",
		Route:       "/",
		UID:         0,
		Command:     "whoami",
	}
	api.apiEPs = append(api.apiEPs, apiRootEP)
	apis = append(apis, api)
	apiPtr := getAPI(api.Name)

	// Create new subrouter
	newApiRouter(apiPtr)

	newHandleFunc(apiPtr, api.apiEPs[0].Name)
	// TODO Test Route var syntax and if it has been used already
	json.NewEncoder(w).Encode(&api)
}

func listAPIs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(apis)
}

func listAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	api := getAPI(apiName)
	json.NewEncoder(w).Encode(api)
}

func listAPIEPs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	epName := vars["endpoint"]
	ep := getApiEP(apiName, epName)
	json.NewEncoder(w).Encode(ep)
}

// newApiRouter
func newApiRouter(api *api) {
	api.router = mainRouter.PathPrefix("/" + api.Name).Subrouter() // "/{apiName}/"
}

func newHandleFunc(a *api, aepName string) {
	a.router.HandleFunc("/"+aepName, generic) // "/{apiName}/{aepName}/"
}

func generic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the generic method executing...")
}

func exApiEP(apiEP *apiEndpoint) error {
	// Executes api based on included methods/commands
	cmd := exec.Command("sleep", "30")
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return err
	}

	return nil
}

func execute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	apiEP := getApiEP(apiName, endpoint)
	// TODO test empty api
	err := exApiEP(apiEP)

	json.NewEncoder(w).Encode(err)
}

func getAPI(name string) *api {
	for i, api := range apis {
		if api.Name == name {
			return &apis[i]
		}
	}
	return nil
}

func getApiEP(_api string, _apiEP string) *apiEndpoint {
	var aep apiEndpoint
	for _, api := range apis {
		if api.Name == _api {
			for _, apiEP := range api.apiEPs {
				if apiEP.Name == _apiEP {
					aep = apiEP
					break
				}
			}
		}
	}
	return &aep
}

func unmarshalApi(r *http.Request) api {
	//TODO error handle and payload validation
	body, _ := ioutil.ReadAll(r.Body)
	var a api
	json.Unmarshal(body, &a)

	return a
}

func unmarshalApiEP(r *http.Request) apiEndpoint {
	//TODO error handle _
	body, _ := ioutil.ReadAll(r.Body)
	var apiEP apiEndpoint
	json.Unmarshal(body, &apiEP)

	return apiEP
}

func startup() {
	// GETs
	mainRouter.HandleFunc("/", listAPIs).Methods("GET")
	mainRouter.HandleFunc("/{api}", listAPI).Methods("GET")
	mainRouter.HandleFunc("/{api}/{endpoint}", listAPIEPs).Methods("GET")

	// POSTs
	mainRouter.HandleFunc("/", createApi).Methods("POST")
	mainRouter.HandleFunc("/{api}", createApiEP).Methods("POST")
	mainRouter.HandleFunc("/{api}/{endpoint}", execute).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", mainRouter))
}

func main() {
	startup()
}
