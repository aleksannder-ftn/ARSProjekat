package services

import (
	"ars_projekat/model"
)

type ConfigurationService struct {
	repo model.ConfigurationRepository
}

func NewConfigurationService(repo model.ConfigurationRepository) ConfigurationService {
	return ConfigurationService{
		repo: repo,
	}
}

func (s ConfigurationService) Add(config *model.Configuration) error {
	err := s.repo.Add(config)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigurationService) Get(name string, version model.Version) (model.Configuration, error) {
	return s.repo.Get(name, version)
}

func (s ConfigurationService) Delete(config model.Configuration) error {
	return s.repo.Delete(config)
}
