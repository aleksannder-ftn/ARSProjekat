package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
)

type ConfigurationGroupService struct {
	repo repositories.IConfigRepository
}

func NewConfigurationGroupService(repo repositories.IConfigRepository) ConfigurationGroupService {
	return ConfigurationGroupService{
		repo: repo,
	}
}

func (s ConfigurationGroupService) Add(configGroup model.ConfigurationGroup) error {
	name := configGroup.Name
	version := model.ToString(configGroup.Version)
	for _, v := range configGroup.Configurations {
		labels := model.SortLabels(v.Labels)
		err := s.repo.AddGroup(name, version, labels, v, context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ConfigurationGroupService) Save(configGroup *model.ConfigurationGroup) error {
	for _, v := range configGroup.Configurations {
		labels := model.SortLabels(v.Labels)
		err := s.repo.AddGroup(configGroup.Name, model.ToString(configGroup.Version), labels, v, context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ConfigurationGroupService) Get(name string, version model.Version, labels string) (*model.ConfigurationGroup, error) {
	return s.repo.GetGroupByParams(name, model.ToString(version), labels)
}

func (s ConfigurationGroupService) Delete(name string, version string, labels string) error {
	return s.repo.DeleteGroupByParams(name, version, labels, context.Background())
}
