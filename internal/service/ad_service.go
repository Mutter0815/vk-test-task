package service

import (
	"github.com/Mutter0815/marketplace/internal/models"
	"github.com/Mutter0815/marketplace/internal/repository"
	"go.uber.org/zap"
)

type AdService struct {
	repo repository.Repository
	log  *zap.SugaredLogger
}

func NewAdService(r repository.Repository, log *zap.SugaredLogger) *AdService {
	return &AdService{repo: r, log: log}
}

func (s *AdService) CreateAd(userID uint, title, desc string, imageURL *string, price uint) (*models.Ad, error) {
	s.log.Debugw("service create ad", "userID", userID, "title", title)
	ad, err := s.repo.CreateAd(userID, title, desc, imageURL, price)
	if err != nil {
		s.log.Errorw("repo create ad failed", "err", err)
	}
	return ad, err
}

func (s *AdService) ListAds(priceMin, priceMax *uint, sortBy, order string, page, pageSize int) ([]*models.Ad, int, error) {
	ads, total, err := s.repo.ListAds(priceMin, priceMax, sortBy, order, page, pageSize)
	if err != nil {
		s.log.Errorw("repo list ads failed", "err", err)
	}
	return ads, total, err
}
