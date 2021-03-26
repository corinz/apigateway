package app

import agw "github.com/corinz/apigateway/pkg/apigateway"

// TODO error check empty or existing name
func (a *app) addSubRouter(api *agw.API) {
	api.Router = a.router.PathPrefix("/" + api.Name).Subrouter() // "/{apiName}/"
}

// TODO error check empty or existing name
func (a *app) newHandleFunc(api *agw.API, subPath string) {
	api.Router.HandleFunc("/"+subPath, Generic) // "/{apiName}/{aepName}/"
}

func (a *app) setupRoutes() {
	// GETs
	a.router.HandleFunc("/", a.ListAPIs).Methods("GET")
	a.router.HandleFunc("/{api}", a.ListAPI).Methods("GET")
	a.router.HandleFunc("/{api}/{endpoint}", a.ListAPIEndpoints).Methods("GET")

	// POSTs
	a.router.HandleFunc("/", a.record(a.CreateAPI)).Methods("POST")
	a.router.HandleFunc("/{api}", a.record(a.CreateAPIEndpoint)).Methods("POST")
	a.router.HandleFunc("/{api}/{endpoint}", a.ExecuteAPIEndpoint).Methods("POST")
}
