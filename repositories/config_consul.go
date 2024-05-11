package repositories

import (
	"ars_projekat/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type ConfigRepository struct {
	cli    *api.Client
	logger *log.Logger
}

func New(logger *log.Logger) (*ConfigRepository, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigRepository{
		cli:    client,
		logger: logger,
	}, nil
}

func (cr *ConfigRepository) GetAll() ([]model.Configuration, error) {
	kv := cr.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	var configurations []model.Configuration
	for _, pair := range data {
		configuration := &model.Configuration{}
		err = json.Unmarshal(pair.Value, configuration)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, *configuration)
	}

	return configurations, nil
}

func (cr *ConfigRepository) GetById(name string, version string) (*model.Configuration, error) {
	kv := cr.cli.KV()
	data, _, err := kv.Get(ConstructConfigKey(name, version), nil)
	if err != nil {
		return nil, err
	}
	configuration := &model.Configuration{}
	err = json.Unmarshal(data.Value, configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func (cr *ConfigRepository) Delete(name string, version string) error {
	kv := cr.cli.KV()

	_, err := kv.Delete(ConstructConfigKey(name, version), nil)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ConfigRepository) Add(config *model.Configuration) (*model.Configuration, error) {
	kv := cr.cli.KV()
	major := strconv.Itoa(config.Version.Major)
	minor := strconv.Itoa(config.Version.Minor)
	patch := strconv.Itoa(config.Version.Patch)
	version := major + "." + minor + "." + patch

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	keyValue := &api.KVPair{Key: ConstructConfigKey(config.Name, version), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}

/*
func getConfigurationsWithLabels(configurations []model.Configuration, requiredLabels []string) []model.Configuration {
	var filteredLabels []model.Configuration

	for _, config := range configurations {
		hasAllLabels := true
		for _, label := range requiredLabels {
			if _, ok := config.Labels[label]; !ok {
				hasAllLabels = false
				break
			}
		}
		if hasAllLabels {
			filteredLabels = append(filteredLabels, config)
		}
	}

	return filteredLabels
} */
