package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

func GenerateToken(user *User, config JWTConfig) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
		Roles:    user.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

func ValidateToken(tokenString string, config JWTConfig) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
}

type RefreshToken struct {
	Token     string
	ExpiresAt time.Time
}
