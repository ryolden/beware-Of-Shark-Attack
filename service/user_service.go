package service

import (
	"api-bewareOfSharkAttack/entity"
	"api-bewareOfSharkAttack/repository"
	"api-bewareOfSharkAttack/utils"
	"errors"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req entity.UserRegisterRequest) (*entity.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	user.Password = "" // Sembunyikan password hashed sebelum di-return
	return user, nil
}

func (s *UserService) Login(req entity.UserLoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("email atau password salah")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
