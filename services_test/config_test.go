// services/configuration_service_test.go
package services

import (
	"ars_projekat/mocks"
	"ars_projekat/model"
	"ars_projekat/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationService_Add(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo.(repositories.ConfigRepository)) // Kastovanje mockRepo

	config := &model.Configuration{Name: "test", Version: "1.0"}
	mockRepo.On("Add", config).Return(config, nil)

	err := service.Add(config)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Add_Error(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo)

	config := &model.Configuration{Name: "test", Version: "1.0"}
	mockRepo.On("Add", config).Return(nil, errors.New("some error"))

	err := service.Add(config)
	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Get(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo)

	config := &model.Configuration{Name: "test", Version: "1.0"}
	mockRepo.On("GetById", "test", "1.0").Return(config, nil)

	result, err := service.Get("test", "1.0")
	assert.NoError(t, err)
	assert.Equal(t, config, result)
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Get_Error(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo)

	mockRepo.On("GetById", "test", "1.0").Return(nil, errors.New("not found"))

	result, err := service.Get("test", "1.0")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Delete(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo)

	config := model.Configuration{Name: "test", Version: "1.0"}
	mockRepo.On("Delete", "test", "1.0").Return(nil)

	err := service.Delete(config)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Delete_Error(t *testing.T) {
	mockRepo := new(mocks.ConfigRepository)
	service := services.NewConfigurationService(mockRepo)

	config := model.Configuration{Name: "test", Version: "1.0"}
	mockRepo.On("Delete", "test", "1.0").Return(errors.New("delete error"))

	err := service.Delete(config)
	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
	mockRepo.AssertExpectations(t)
}
