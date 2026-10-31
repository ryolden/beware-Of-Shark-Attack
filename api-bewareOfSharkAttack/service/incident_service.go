package service

import (
	"api-bewareOfSharkAttack/entity"
	"api-bewareOfSharkAttack/repository"
)

type IncidentService struct {
	repo *repository.IncidentRepository
}

func NewIncidentService(repo *repository.IncidentRepository) *IncidentService {
	return &IncidentService{repo: repo}
}

func (s *IncidentService) GetAll(limit, page int) ([]entity.Incident, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.repo.GetAll(limit, offset)
}

func (s *IncidentService) GetByID(id int) (*entity.Incident, error) {
	return s.repo.GetByID(id)
}

func (s *IncidentService) Create(inc *entity.Incident) error {
	return s.repo.Create(inc)
}

func (s *IncidentService) Delete(id int) error {
	return s.repo.Delete(id)
}
