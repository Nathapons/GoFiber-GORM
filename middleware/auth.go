package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	jwtSecretKey := "1234"
	_, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	fmt.Println("err======>", err)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
