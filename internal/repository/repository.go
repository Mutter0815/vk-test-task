package repository

import "github.com/Mutter0815/marketplace/internal/models"

type Repository interface {
	CreateUser(username, password string) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)
	CreateAd(userID uint, title, desc string, imageURL *string, price uint) (*models.Ad, error)
	ListAds(priceMin, priceMax *uint, sortBy, order string, page, pageSize int) ([]*models.Ad, int, error)
}
