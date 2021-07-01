package plugin

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers stores the list of http.HandlerFunc functions for the different resource calls
type Handlers struct {
	Projects     http.HandlerFunc
	Clusters http.HandlerFunc
	Mongos http.HandlerFunc
	Disks http.HandlerFunc
	Databases http.HandlerFunc
}

// GetRouter creates the gorilla/mux router for the HTTP handlers
func GetRouter(h Handlers) *mux.Router {
	router := mux.NewRouter()
	router.Path("/projects").Methods("GET").HandlerFunc(h.Projects)
	router.Path("/clusters").Methods("GET").HandlerFunc(h.Clusters)
	router.Path("/mongos").Methods("GET").HandlerFunc(h.Mongos)
	router.Path("/disks").Methods("GET").HandlerFunc(h.Disks)
	router.Path("/databases").Methods("GET").HandlerFunc(h.Databases)

	return router
}
