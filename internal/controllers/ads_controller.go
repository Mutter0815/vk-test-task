package controllers

import (
	"net/http"

	"github.com/Mutter0815/marketplace/internal/dto"
	"github.com/Mutter0815/marketplace/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ADSController struct {
	service *service.AdService
	log     *zap.SugaredLogger
}

func NewADSController(s *service.AdService, log *zap.SugaredLogger) *ADSController {
	return &ADSController{service: s, log: log}
}

func (ac *ADSController) CreateAd(c *gin.Context) {
	userID := c.GetUint("userID")
	var req dto.CreateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.log.Warnw("create ad bind failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("creating ad", "userID", userID, "title", req.Title)
	ad, err := ac.service.CreateAd(userID, req.Title, req.Description, req.ImageURL, req.Price)
	if err != nil {
		ac.log.Errorw("create ad failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("ad created", "id", ad.ID)
	c.JSON(http.StatusCreated, dto.AdResponseFromModel(ad, userID))

}

func (ac *ADSController) ListAds(c *gin.Context) {

	var q dto.ListAdsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		ac.log.Warnw("list ads bind failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("listing ads", "page", q.Page, "page_size", q.PageSize)

	ads, total, err := ac.service.ListAds(
		q.PriceMin, q.PriceMax,
		q.SortBy, q.Order,
		q.Page, q.PageSize,
	)
	if err != nil {
		ac.log.Errorw("list ads failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userID uint
	if uid, ok := c.Get("userID"); ok {
		userID = uid.(uint)
	}

	resp := dto.ListAdsResponse{
		Ads:      dto.AdListFromModels(ads, userID),
		Page:     q.Page,
		PageSize: q.PageSize,
		Total:    total,
	}
	c.JSON(http.StatusOK, resp)
}
