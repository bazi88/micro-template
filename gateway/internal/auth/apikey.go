package auth

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

type APIKeyConfig struct {
	Header string
	Secret string
}

func ValidateAPIKey(c *fiber.Ctx, config APIKeyConfig) error {
	apiKey := c.Get(config.Header)
	if apiKey == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing API key")
	}

	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(config.Secret)) != 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid API key")
	}

	return nil
}
