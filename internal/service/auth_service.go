package service

import (
	"github.com/Mutter0815/marketplace/internal/models"
	"github.com/Mutter0815/marketplace/internal/repository"
	"go.uber.org/zap"
)

// AuthService provides user related operations.
type AuthService struct {
	repo repository.Repository
	log  *zap.SugaredLogger
}

func NewAuthService(r repository.Repository, log *zap.SugaredLogger) *AuthService {
	return &AuthService{repo: r, log: log}
}

func (s *AuthService) Register(username, password string) (*models.User, error) {
	s.log.Debugw("service register", "username", username)
	u, err := s.repo.CreateUser(username, password)
	if err != nil {
		s.log.Errorw("create user failed", "err", err)
	}
	return u, err
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	s.log.Debugw("service login", "username", username)
	u, err := s.repo.AuthenticateUser(username, password)
	if err != nil {
		s.log.Warnw("authenticate failed", "err", err)
	}
	return u, err
}
