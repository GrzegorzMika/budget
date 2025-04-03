package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const jwksURL = `http://127.0.0.1:8080/realms/budget/protocol/openid-connect/certs`

type AuthorizationHeader struct {
	Bearer string `header:"Authorization"`
}

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := new(AuthorizationHeader)
		if err := c.Bind().Header(authHeader); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		header := strings.Fields(authHeader.Bearer)
		if len(header) != 2 || header[0] != "Bearer" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		token := header[1]
		jwksKeySet, err := jwk.Fetch(c.Context(), jwksURL)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		_, err = jwt.Parse([]byte(token), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
