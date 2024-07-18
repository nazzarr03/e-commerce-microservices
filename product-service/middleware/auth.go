package middleware

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("JWT_SECRET")

type JWTClaims struct {
	UserID uint `json:"id"`
	jwt.StandardClaims
}

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing JWT",
			})
		}

		token := authHeader[len("Bearer "):]

		claims, err := ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

func ValidateToken(token string) (*JWTClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid JWT claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
