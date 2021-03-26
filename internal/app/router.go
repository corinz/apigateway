package app

import agw "github.com/corinz/apigateway/pkg/apigateway"

// TODO error check empty or existing name
func (a *app) addSubRouter(api *agw.API) {
	api.Router = a.router.PathPrefix("/" + api.Name).Subrouter() // "/{apiName}/"
}

// TODO error check empty or existing name
func (a *app) newHandleFunc(api *agw.API, subPath string) {
	api.Router.HandleFunc("/"+subPath, generic) // "/{apiName}/{aepName}/"
}

func (a *app) setupRoutes() {
	// GETs
	a.router.HandleFunc("/", a.listAPIs).Methods("GET")
	a.router.HandleFunc("/{api}", a.listAPI).Methods("GET")
	a.router.HandleFunc("/{api}/{endpoint}", a.listAPIEndpoints).Methods("GET")

	// POSTs
	a.router.HandleFunc("/", a.record(a.createAPI)).Methods("POST")
	a.router.HandleFunc("/{api}", a.record(a.createAPIEndpoint)).Methods("POST")
	a.router.HandleFunc("/{api}/{endpoint}", a.executeAPIEndpoint).Methods("POST")
}
