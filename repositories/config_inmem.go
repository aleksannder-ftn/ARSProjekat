package repositories

import "ars_projekat/model"

type ConfigInMemoryRepository struct {
	configs map[string]model.Configuration
}

func (c ConfigInMemoryRepository) Get() {
	panic("implement me")
}

func NewConfigInMemoryRepository() model.ConfigurationRepository {
	return ConfigInMemoryRepository{
		configs: make(map[string]model.Configuration),
	}
}
