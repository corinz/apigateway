package main

import (
	"encoding/json"
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

	appendApiEP(apiName, apiEP)

	// TODO Add this new EP to the the API's new router (child router)
	json.NewEncoder(w).Encode(apiEP)
}

func appendApiEP(apiName string, apiEP apiEndpoint) {
	for i, api := range apis {
		if api.Name == apiName {
			// Appends to global var
			apis[i].apiEPs = append(apis[i].apiEPs, apiEP)
			break
		}
	}
}

func createApi(w http.ResponseWriter, r *http.Request) {
	api := unmarshalApi(r)
	//TODO Test to see if API already exists and validate the Route is unique/valid

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

	// Create new router
	newApiRouter(api.Name)

	newHandleFunc(api) // routes to a placeholder method
	newHandle(api)     // adds child router to main router
	// TODO Test Route var syntax and if it has been used already
	//r.HandleFunc("/", exApiEP(apiEP)) // needs updated method
	//mainRouter.Handle(route, r)

	json.NewEncoder(w).Encode(api)
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
func newApiRouter(apiName string) {
	api := getAPI(apiName)
	api.router = newRouter()
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	return r
}

func newHandleFunc(a api) {
	a.router.HandleFunc("/", generic)
}

func newHandle(a api) {
	mainRouter.Handle(a.Name, a.router)
}

func generic(w http.ResponseWriter, r *http.Request) {

}

func exApiEP(apiEP apiEndpoint) func(w http.ResponseWriter, r *http.Request) {
	// Executes api based on included methods/commands
	cmd := exec.Command("sleep", "1")
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(apiEP)
	}
}

func execute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["api"]
	endpoint := vars["endpoint"]

	apiEP := getApiEP(apiName, endpoint)
	// TODO test empty api
	exApiEP(apiEP)
}

func getAPI(name string) api {
	var a api
	for i, api := range apis {
		if api.Name == name {
			a = apis[i]
			break
		}
	}
	return a
}

func getApiEP(_api string, _apiEP string) apiEndpoint {
	var a apiEndpoint
	for _, api := range apis {
		if api.Name == _api {
			for _, apiEP := range api.apiEPs {
				if apiEP.Name == _apiEP {
					a = apiEP
					break
				}
			}
		}
	}
	return a
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
