package handlers

import (
	"ars_projekat/model"
	"ars_projekat/services"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"strings"
)

type ConfigurationGroupHandler struct {
	service services.ConfigurationGroupService
}

func NewConfigurationGroupHandler(service services.ConfigurationGroupService) ConfigurationGroupHandler {
	return ConfigurationGroupHandler{service: service}
}

func (cg ConfigurationGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := convertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.service.Get(name, versionModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, cGroup)
}

func (cg ConfigurationGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
	cType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(cType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cfgGroup, err := decodeGroupBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	check, err := cg.service.Get(cfgGroup.Name, cfgGroup.Version)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if check.Id != 0 {
		err := errors.New("config already exists")
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	cg.service.Add(*cfgGroup)
	renderJSON(w, cfgGroup)
}

func (cg ConfigurationGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		err := errors.New("expect application/json content-type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cfg, err := decodeGroupBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCfg, err := cg.service.Update(*cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, updatedCfg)
}

func (cg ConfigurationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := convertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.service.Get(name, versionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := cg.service.Delete(cGroup)

	if ok != nil {
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, "successfully deleted")
}

func decodeGroupBody(r io.Reader) (*model.ConfigurationGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cg model.ConfigurationGroup
	if err := dec.Decode(&cg); err != nil {
		return nil, err
	}

	return &cg, nil
}
