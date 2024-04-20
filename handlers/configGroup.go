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
	groupService  services.ConfigurationGroupService
	configService services.ConfigurationService
}

func NewConfigurationGroupHandler(groupService services.ConfigurationGroupService, configService services.ConfigurationService) ConfigurationGroupHandler {
	return ConfigurationGroupHandler{groupService: groupService, configService: configService}
}

func (cg ConfigurationGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := convertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.groupService.Get(name, versionModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	for _, cfg := range cfgGroup.Configurations {
		check, err := cg.configService.Get(cfg.Name, cfg.Version)
		if err != nil && !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if check.Id != 0 {
			continue
		}

		err = cg.configService.Add(&cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	check, err := cg.groupService.Get(cfgGroup.Name, cfgGroup.Version)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if check.Id != 0 {
		err := errors.New("config already exists")
		http.Error(w, err.Error(), http.StatusConflict)
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
	versionModel, err := convertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.groupService.Get(name, versionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := cg.groupService.Delete(cGroup)

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
