package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
)

type ConfigurationService struct {
	repo repositories.ConfigRepository
}

func NewConfigurationService(repo repositories.ConfigRepository) ConfigurationService {
	return ConfigurationService{
		repo: repo,
	}
}

func (s ConfigurationService) Add(config *model.Configuration) error {
	_, err := s.repo.Add(config)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigurationService) Get(name string, version string) (*model.Configuration, error) {
	return s.repo.GetById(name, version)
}

func (s ConfigurationService) Delete(config model.Configuration) error {
	ver := model.ToString(config.Version)
	return s.repo.Delete(config.Name, ver)
}
