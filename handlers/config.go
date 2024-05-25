package handlers

import (
	"ars_projekat/model"
	"ars_projekat/services"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

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

// swagger:route GET /configs/{name}/{version} configuration getConfiguration
// Get configuration by name and version
//
// responses:
//
//	404: ErrorResponse
//	200: Configuration
func (c ConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	config, err := c.service.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, config)
}

// swagger:route POST /configs configuration upsertConfiguration
// Add or update a configuration
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	409: ErrorResponse
//	201: Configuration
func (c ConfigurationHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cfg, err := decodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ver := model.ToString(cfg.Version)
	check, err := c.service.Get(cfg.Name, ver)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if check != nil {
		err := errors.New("config already exists")
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = c.service.Add(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(w, cfg)
}

// swagger:route DELETE /configs/{name}/{version} configuration deleteConfiguration
// Delete a configuration by name and version
//
// responses:
//
//	404: ErrorResponse
//	204: NoContent
func (c ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	config, err := c.service.Get(name, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := c.service.Delete(*config)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeBody(r io.Reader) (*model.Configuration, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var configuration model.Configuration
	if err := dec.Decode(&configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	marshal, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(marshal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
