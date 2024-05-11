package handlers

import (
	"ars_projekat/model"
	"ars_projekat/services"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strconv"
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

// Get  /configs/{name}/{version}
func (c ConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := ConvertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, config)
}

// Post /configs/
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

	check, err := c.service.Get(cfg.Name, cfg.Version)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if check.Id != 0 {
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

// Delete /configs/{name}/{version}
func (c ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := ConvertVersion(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := c.service.Delete(config)

	if ok != nil {
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, "successfully deleted")
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

func ConvertVersion(version string) (model.Version, error) {
	split := strings.Split(version, ".")
	if len(split) != 3 {
		return model.Version{}, errors.New("version incorrect")
	}

	major, err := strconv.Atoi(split[0])
	if err != nil {
		return model.Version{}, errors.New("failed converting")
	}
	minor, err := strconv.Atoi(split[1])
	if err != nil {
		return model.Version{}, errors.New("failed converting")
	}
	patch, err := strconv.Atoi(split[2])
	if err != nil {
		return model.Version{}, errors.New("failed converting")
	}

	versionModel := model.Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}

	return versionModel, nil
}
