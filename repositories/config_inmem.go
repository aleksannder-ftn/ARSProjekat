package repositories

import (
	"ars_projekat/model"
	"errors"
	"fmt"
)

type ConfigInMemoryRepository struct {
	configs map[string]model.Configuration
}

func (c ConfigInMemoryRepository) Add(config model.Configuration) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

func (c ConfigInMemoryRepository) Get(name string, version int) (model.Configuration, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Configuration{}, errors.New("config not found")
	}
	return config, nil
}

func (c ConfigInMemoryRepository) Delete(config model.Configuration) error {
	key := fmt.Sprintf("%s%d", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}

func (c ConfigInMemoryRepository) Update(config model.Configuration) error {
	key := fmt.Sprintf("%s%d", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("config not found")
	}
	c.configs[key] = config
	return nil
}

func NewConfigInMemoryRepository() model.ConfigurationRepository {
	return ConfigInMemoryRepository{
		configs: make(map[string]model.Configuration),
	}
}
