package app

import (
	"net/http"
)

func (a *app) setupRoutes() {
	// swagger docs
	docs := http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/swagger-ui/dist")))
	a.router.PathPrefix("/docs/").Handler(docs)
	a.router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", 301)
	})

	// reroute homepage to docs
	a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", 301)
	})

	// API Routes ../api/{api}/{endpoint}
	// GETs
	a.router.HandleFunc("/api", a.listAPIs).Methods("GET")
	a.router.HandleFunc("/api/{api}", a.listAPI).Methods("GET")
	a.router.HandleFunc("/api/{api}/{endpoint}", a.listAPIEndpoints).Methods("GET")

	// POSTs
	a.router.HandleFunc("/api", a.record(a.createAPI)).Methods("POST")
	a.router.HandleFunc("/api/{api}", a.record(a.createAPIEndpoint)).Methods("POST")

	// DELETEs
	a.router.HandleFunc("/api/{api}", a.record(a.delete)).Methods("DELETE")
	a.router.HandleFunc("/api/{api}/{endpoint}", a.record(a.delete)).Methods("DELETE")

	// User-created endpoints ../{api}/{endpoint}
	// GETs
	a.router.HandleFunc("/{api}", a.executeAPI).Methods("GET")
	a.router.HandleFunc("/{api}/{endpoint}", a.executeAPIEndpoint).Methods("GET")
}
