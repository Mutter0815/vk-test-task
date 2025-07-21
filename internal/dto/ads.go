package dto

import (
	"time"

	"github.com/Mutter0815/marketplace/internal/models"
)

type CreateAdRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"required,min=1,max=1000"`
	ImageURL    *string `json:"image_url,omitempty" binding:"omitempty,url,max=255"`
	Price       uint    `json:"price" binding:"required,gt=0"`
}

type AdResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Author      string  `json:"author"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url,omitempty"`
	Price       uint    `json:"price"`
	CreatedAt   string  `json:"created_at"`
	IsMine      bool    `json:"isMine,omitempty"`
}

func AdResponseFromModel(ad *models.Ad, currentUserID uint) AdResponse {
	return AdResponse{
		ID:          ad.ID,
		UserID:      ad.UserID,
		Author:      ad.Author,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		CreatedAt:   ad.CreatedAt.Format(time.RFC3339),
		IsMine:      ad.UserID == currentUserID,
	}
}

type ListAdsResponse struct {
	Ads      []AdResponse `json:"ads"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Total    int          `json:"total"`
}

type ListAdsQuery struct {
	Page     int    `form:"page,default=1" binding:"gte=1"`
	PageSize int    `form:"page_size,default=10" binding:"gte=1,lte=100"`
	PriceMin *uint  `form:"price_min"`
	PriceMax *uint  `form:"price_max"`
	SortBy   string `form:"sort_by,default=date" binding:"oneof=date price"`
	Order    string `form:"order,default=desc" binding:"oneof=asc desc"`
}

func AdListFromModels(ads []*models.Ad, currentUserID uint) []AdResponse {
	out := make([]AdResponse, len(ads))
	for i, a := range ads {
		out[i] = AdResponseFromModel(a, currentUserID)
	}
	return out
}
