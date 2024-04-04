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

func (s ConfigurationGroupService) Add(configGroup model.ConfigurationGroup) {
	s.repo.Add(configGroup)
}

func (s ConfigurationGroupService) Get(name string, version int) (model.ConfigurationGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigurationGroupService) Update(configGroup model.ConfigurationGroup) (model.ConfigurationGroup, error) {
	return s.repo.Update(configGroup)
}

func (s ConfigurationGroupService) Delete(configGroup model.ConfigurationGroup) (model.ConfigurationGroup, error) {
	return s.repo.Delete(configGroup)
}
