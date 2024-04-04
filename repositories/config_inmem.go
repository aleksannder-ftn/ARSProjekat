package repositories

import (
	"ars_projekat/model"
	"errors"
	"fmt"
)

type ConfigInMemoryRepository struct {
	configs map[string]model.Configuration
}

// Add implements model.ConfigurationRepository.
func (c ConfigInMemoryRepository) Add(config model.Configuration) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

// Get implements model.ConfigurationRepository.
func (c ConfigInMemoryRepository) Get(name string, version int) (model.Configuration, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Configuration{}, errors.New("config not found")
	}
	println(config)
	return model.Configuration{}, errors.New("config not found")
}

func NewConfigInMemoryRepository() model.ConfigurationRepository {
	return ConfigInMemoryRepository{
		configs: make(map[string]model.Configuration),
	}
}
