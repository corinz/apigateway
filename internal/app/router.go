package app

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
