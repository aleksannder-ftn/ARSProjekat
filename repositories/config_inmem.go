package repositories

import (
	"ars_projekat/model"
	"errors"
	"fmt"
)

type ConfigInMemoryRepository struct {
	configs map[string]model.Configuration
}

func (c ConfigInMemoryRepository) Add(config *model.Configuration) {
	key := fmt.Sprintf("%s/%#v", config.Name, config.Version)
	c.configs[key] = *config
}

func (c ConfigInMemoryRepository) Get(name string, version model.Version) (model.Configuration, error) {
	key := fmt.Sprintf("%s/%#v", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Configuration{}, errors.New("config not found")
	}
	return config, nil
}

func (c ConfigInMemoryRepository) Delete(config model.Configuration) error {
	key := fmt.Sprintf("%s/%#v", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}

func (c ConfigInMemoryRepository) Update(config model.Configuration) (model.Configuration, error) {
	key := fmt.Sprintf("%s/%#v", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return model.Configuration{}, errors.New("config not found")
	}
	config.Version = model.IncrementVersion(config.Version)
	c.Add(&config)
	return config, nil
}

func NewConfigInMemoryRepository() model.ConfigurationRepository {
	return ConfigInMemoryRepository{
		configs: make(map[string]model.Configuration),
	}
}
