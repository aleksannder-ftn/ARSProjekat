package services

import (
	"ars_projekat/model"
)

type ConfigurationGroupService struct {
	repo model.ConfigurationGroupRepository
}

func NewConfigurationGroupService(repo model.ConfigurationGroupRepository) ConfigurationGroupService {
	return ConfigurationGroupService{
		repo: repo,
	}
}

func (s ConfigurationGroupService) Add(configGroup model.ConfigurationGroup) error {
	err := s.repo.Add(&configGroup)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigurationGroupService) Get(name string, version model.Version) (model.ConfigurationGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigurationGroupService) Delete(configGroup model.ConfigurationGroup) error {
	return s.repo.Delete(configGroup)
}
