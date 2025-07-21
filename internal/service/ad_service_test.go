package service

import (
	"testing"

	"github.com/Mutter0815/marketplace/internal/models"
	"go.uber.org/zap"
)

type adRepoMock struct {
	createAdCalled bool
	createAdArgs   struct {
		uid         uint
		title, desc string
		img         *string
		price       uint
	}
	createAdResp *models.Ad
	createAdErr  error
}

func (m *adRepoMock) CreateUser(u, p string) (*models.User, error)       { return nil, nil }
func (m *adRepoMock) AuthenticateUser(u, p string) (*models.User, error) { return nil, nil }
func (m *adRepoMock) CreateAd(uid uint, title, desc string, img *string, price uint) (*models.Ad, error) {
	m.createAdCalled = true
	m.createAdArgs.uid = uid
	m.createAdArgs.title = title
	m.createAdArgs.desc = desc
	m.createAdArgs.img = img
	m.createAdArgs.price = price
	return m.createAdResp, m.createAdErr
}
func (m *adRepoMock) ListAds(min, max *uint, sortBy, order string, page, pageSize int) ([]*models.Ad, int, error) {
	return nil, 0, nil
}

func TestAdServiceCreateAd(t *testing.T) {
	repo := &adRepoMock{createAdResp: &models.Ad{ID: 2, Title: "t"}}
	svc := NewAdService(repo, zap.NewNop().Sugar())
	ad, err := svc.CreateAd(1, "t", "d", nil, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.createAdCalled {
		t.Fatal("CreateAd not called")
	}
	if repo.createAdArgs.title != "t" || repo.createAdArgs.price != 10 {
		t.Fatal("wrong args")
	}
	if ad.ID != 2 {
		t.Fatal("wrong result")
	}
}
func TestAdServiceCreateAdWithImageURL(t *testing.T) {
	repo := &adRepoMock{createAdResp: &models.Ad{ID: 3}}
	svc := NewAdService(repo, zap.NewNop().Sugar())
	url := "http://example.com/img.png"
	ad, err := svc.CreateAd(2, "title", "desc", &url, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.createAdArgs.img == nil || *repo.createAdArgs.img != url {
		t.Fatalf("image url not passed")
	}
	if ad.ID != 3 {
		t.Fatal("wrong result")
	}
}
