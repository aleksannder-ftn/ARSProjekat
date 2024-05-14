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

type ConfigurationGroupHandler struct {
	groupService services.ConfigurationGroupService
}

func NewConfigurationGroupHandler(groupService services.ConfigurationGroupService) ConfigurationGroupHandler {
	return ConfigurationGroupHandler{groupService: groupService}
}

func (cg ConfigurationGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	labels := strings.Split(mux.Vars(r)["labels"], ";")

	var labelString string
	for i, v := range labels {
		if i == len(labels)-1 {
			labelString += v
		} else {
			labelString += v
			labelString += ";"
		}
	}
	versionModel, err := model.ToVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.groupService.Get(name, *versionModel, labelString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, cGroup)
}

func (cg ConfigurationGroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := model.ToVersion(version)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.groupService.Get(name, *versionModel, "")

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

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

	config, err := decodeBody(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range cGroup.Configurations {
		if v.Name == config.Name && v.Version == config.Version {
			err := errors.New("config is already added")
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}

	cGroup.Configurations = append(cGroup.Configurations, *config)

	err = cg.groupService.Save(cGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, cGroup)
}

func (cg ConfigurationGroupHandler) Upsert(w http.ResponseWriter, r *http.Request) {
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

	err = cg.groupService.Add(*cfgGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, cfgGroup)
}

func (cg ConfigurationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	labels := strings.Split(mux.Vars(r)["labels"], ";")

	var labelString string
	for i, v := range labels {
		if i == len(labels)-1 {
			labelString += v
		} else {
			labelString += v
			labelString += ";"
		}
	}
	versionModel, err := model.ToVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	check, err := cg.groupService.Get(name, *versionModel, labelString)
	if check == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := cg.groupService.Delete(name, version, labelString)

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
