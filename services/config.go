package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
)

type ConfigurationService struct {
	repo repositories.IConfigRepository
}

// NewConfigurationService creates a new instance of ConfigurationService.
func NewConfigurationService(repo repositories.IConfigRepository) ConfigurationService {
	return ConfigurationService{
		repo: repo,
	}
}

// Add adds a new configuration.
func (s ConfigurationService) Add(config *model.Configuration, ctx context.Context) error {
	_, err := s.repo.Add(config, ctx)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a configuration by name and version.
func (s ConfigurationService) Get(name string, version string, ctx context.Context) (*model.Configuration, error) {
	return s.repo.GetById(name, version, ctx)
}

// Delete deletes a configuration.
func (s ConfigurationService) Delete(config model.Configuration, ctx context.Context) error {
	ver := model.ToString(config.Version)
	return s.repo.Delete(config.Name, ver, ctx)
}
