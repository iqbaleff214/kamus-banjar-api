package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Auth validates the Bearer JWT and sets "userID" and "userRole" in Locals.
func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid Authorization header")
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			return fiber.NewError(fiber.StatusUnauthorized, "token expired")
		}

		sub, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)
		if sub == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token subject")
		}

		c.Locals("userID", sub)
		c.Locals("userRole", role)
		return c.Next()
	}
}

// Role enforces that the authenticated user has the required role.
// Must be chained after Auth.
func Role(required string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("userRole").(string)
		if role != required {
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
		}
		return c.Next()
	}
}
