package mocks

import (
	"ars_projekat/model"

	"github.com/stretchr/testify/mock"
)

type ConfigRepositoryInterface interface {
    Add(config *model.Configuration) (*model.Configuration, error)
    GetById(name string, version string) (*model.Configuration, error)
    Delete(name string, version string) error
}

type ConfigRepository struct {
	mock.Mock
}

func (m *ConfigRepository) Add(config *model.Configuration) (*model.Configuration, error) {
	args := m.Called(config)
	return args.Get(0).(*model.Configuration), args.Error(1)
}

func (m *ConfigRepository) GetById(name string, version string) (*model.Configuration, error) {
	args := m.Called(name, version)
	return args.Get(0).(*model.Configuration), args.Error(1)
}

func (m *ConfigRepository) Delete(name string, version string) error {
	args := m.Called(name, version)
	return args.Error(0)
}
