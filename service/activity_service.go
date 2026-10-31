package service

import (
	"api-bewareOfSharkAttack/entity"
	"api-bewareOfSharkAttack/repository"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) GetAll() ([]entity.Activity, error) {
	return s.repo.GetAll()
}

func (s *ActivityService) GetByID(id int) (*entity.Activity, error) {
	return s.repo.GetByID(id)
}

func (s *ActivityService) Create(act *entity.Activity) error {
	return s.repo.Create(act)
}

func (s *ActivityService) Update(id int, act *entity.Activity) error {
	return s.repo.Update(id, act)
}

func (s *ActivityService) Delete(id int) error {
	return s.repo.Delete(id)
}
