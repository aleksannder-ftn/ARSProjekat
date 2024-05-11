package repositories

import (
	"ars_projekat/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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

func (pr *ConfigRepository) GetAll() ([]model.Configuration, error) {
	kv := pr.cli.KV()
	data, _, err := kv.List(allConfigs, nil)
	if err != nil {
		return nil, err
	}

	configurations := []model.Configuration{}
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
	if data == nil {
		return nil, errors.New("not found")
	}
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
	version := model.ToString(config.Version)

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

// Config group
func (cr *ConfigRepository) GetAllGroups() ([]model.ConfigurationGroup, error) {
	kv := cr.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {
		return nil, err
	}

	var groups []model.ConfigurationGroup
	for _, pair := range data {
		cg := &model.ConfigurationGroup{}
		err = json.Unmarshal(pair.Value, cg)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *cg)
	}

	return groups, nil
}

func (cr *ConfigRepository) GetGroupByParams(name string, version string, labels string) (*model.ConfigurationGroup, error) {
	kv := cr.cli.KV()

	var key string
	if len(labels) == 0 {
		key = ConstructConfigGroupKey(name, version, "", "")
	} else {
		key = ConstructConfigGroupKey(name, version, labels, "")
	}

	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, err
	}
	cg := &model.ConfigurationGroup{}
	cg.Name = name
	ver, err := model.ToVersion(version)
	if err != nil {
		return nil, err
	}
	cg.Version = *ver

	for _, pair := range data {
		config := &model.Configuration{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		cg.Configurations = append(cg.Configurations, *config)
	}

	return cg, nil
}

func (cr *ConfigRepository) AddGroup(name string, version string, labels string, configs model.Configuration) error {
	kv := cr.cli.KV()

	data, err := json.Marshal(configs)
	if err != nil {
		return err
	}

	keyValue := &api.KVPair{Key: ConstructConfigGroupKey(name, version, labels, configs.Name), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ConfigRepository) DeleteGroupById(name string, version string) error {
	kv := cr.cli.KV()

	_, err := kv.Delete(ConstructConfigKey(name, version), nil)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ConfigRepository) DeleteGroupByParams(name string, version string, labels string, configs model.Configuration) error {
	kv := cr.cli.KV()

	var key string
	if len(labels) == 0 {
		key = ConstructConfigGroupKey(name, version, "", "")
	} else {
		key = ConstructConfigGroupKey(name, version, labels, "")
	}
	_, err := kv.DeleteTree(key, nil)
	if err != nil {
		return err
	}

	return nil
}
