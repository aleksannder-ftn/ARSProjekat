package handlers

import (
	"ars_projekat/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigurationHandler struct {
	service services.ConfigurationService
}

func NewConfigurationHandler(service services.ConfigurationService) ConfigurationHandler {
	return ConfigurationHandler{
		service: service,
	}
}

// GET /configurations/{name}/{version}
func (c ConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	config, err := c.service.Get(name, versionInt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	println(config)
	// ovde u json prebaciti
}
