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

	cg.service.Add(*cfgGroup)

	renderJSON(w, cfgGroup)
}

func (cg ConfigurationGroupHandler) Update(w http.ResponseWriter, r *http.Request) {

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
