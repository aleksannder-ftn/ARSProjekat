package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
)

type IdempotencyService struct {
	repo repositories.ConfigRepository
}

func NewIdempotencyService(repo repositories.ConfigRepository) IdempotencyService {
	return IdempotencyService{
		repo: repo,
	}
}

func (i IdempotencyService) Add(req *model.IdempotencyRequest) error {
	_, err := i.repo.AddIdempotencyRequest(req, context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (i IdempotencyService) Get(key string) (bool, error) {
	exists, err := i.repo.GetIdempotencyRequestByKey(key, context.Background())
	if err != nil {
		return false, err
	}

	return exists, nil
}
