package repositories

import (
	"ars_projekat/model"
	"encoding/json"
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
	data, _, err := kv.List(all, nil)
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
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}
