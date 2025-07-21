package service

import (
	"testing"

	"github.com/Mutter0815/marketplace/internal/models"
	"go.uber.org/zap"
)

type authRepoMock struct {
	createUserCalled bool
	createUserArgs   struct{ u, p string }
	createUserResp   *models.User
	createUserErr    error
}

func (m *authRepoMock) CreateUser(u, p string) (*models.User, error) {
	m.createUserCalled = true
	m.createUserArgs.u = u
	m.createUserArgs.p = p
	return m.createUserResp, m.createUserErr
}
func (m *authRepoMock) AuthenticateUser(u, p string) (*models.User, error) { return nil, nil }
func (m *authRepoMock) CreateAd(uid uint, t, d string, img *string, price uint) (*models.Ad, error) {
	return nil, nil
}
func (m *authRepoMock) ListAds(min, max *uint, sortBy, order string, page, pageSize int) ([]*models.Ad, int, error) {
	return nil, 0, nil
}

func TestAuthServiceRegister(t *testing.T) {
	repo := &authRepoMock{createUserResp: &models.User{ID: 1, Username: "test"}}
	svc := NewAuthService(repo, zap.NewNop().Sugar())
	u, err := svc.Register("test", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.createUserCalled {
		t.Fatal("CreateUser not called")
	}
	if repo.createUserArgs.u != "test" || repo.createUserArgs.p != "pass" {
		t.Fatal("wrong args")
	}
	if u.Username != "test" {
		t.Fatal("unexpected user")
	}
}
