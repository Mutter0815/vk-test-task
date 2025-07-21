package controllers

import (
	"net/http"
	"time"

	"github.com/Mutter0815/marketplace/internal/dto"
	"github.com/Mutter0815/marketplace/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type AuthController struct {
	service   *service.AuthService
	jwtSecret []byte
	log       *zap.SugaredLogger
}

func NewAuthController(s *service.AuthService, secret []byte, log *zap.SugaredLogger) *AuthController {
	return &AuthController{service: s, jwtSecret: secret, log: log}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.log.Warnw("register request bind failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("registering user", "username", req.Username)
	user, err := ac.service.Register(req.Username, req.Password)
	if err != nil {
		ac.log.Errorw("register failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("user registered", "id", user.ID)
	resp := dto.RegisterResponse{ID: user.ID, Username: user.Username, Message: "Пользоватен зарегистрирован"}
	c.JSON(http.StatusOK, resp)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.log.Warnw("login request bind failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ac.log.Infow("user login", "username", req.Username)
	user, err := ac.service.Login(req.Username, req.Password)
	if err != nil {
		ac.log.Warnw("login failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	ts, err := token.SignedString(ac.jwtSecret)
	if err != nil {
		ac.log.Errorw("sign token failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "невозможно добавить токен"})
		return
	}
	ac.log.Infow("token issued", "userID", user.ID)
	resp := dto.LoginResponse{Token: ts}
	c.JSON(http.StatusOK, resp)
}
