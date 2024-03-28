package services

import (
	"ars_projekat/model"
	"fmt"
)

type ConfigurationService struct {
	repo model.ConfigurationRepository
}

func NewConfigurationService(repo model.ConfigurationRepository) ConfigurationService {
	return ConfigurationService{
		repo: repo,
	}
}

// TODO add CRUD
func (s ConfigurationService) Hello() {
	fmt.Println("Hello from config service")
}
