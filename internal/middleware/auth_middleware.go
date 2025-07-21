package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func JWTAuthMiddleware(secret []byte, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn("invalid auth header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be 'Bearer <token>'"})
			c.Abort()
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			log.Warnw("invalid token", "err", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["sub"].(float64))
		log.Debugw("token parsed", "userID", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
