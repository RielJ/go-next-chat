package api

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/rielj/go-next-chat/internal/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var authorizationHeader string
			authorizationHeader = c.Request().Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				authorizationCookie, err := c.Cookie(authorizationHeaderKey)
				if err != nil {
					return echo.ErrUnauthorized
				}
				authorizationHeader = authorizationCookie.Value
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) != 2 {
				return echo.ErrUnauthorized
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				return echo.ErrUnauthorized
			}

			authorizationToken := fields[1]
			payload, err := tokenMaker.VerifyToken(authorizationToken)
			if err != nil {
				return echo.ErrUnauthorized
			}

			// cookies := new(http.Cookie)
			// cookies.Name = authorizationHeaderKey
			// cookies.Value = authorizationHeader
			// cookies.Path = "/"
			// cookies.HttpOnly = true
			// cookies.Secure = true

			// c.SetCookie(cookies)
			// c.Response().Header().Set(authorizationPayloadKey, payload)
			c.Set(authorizationPayloadKey, payload)
			return next(c)
		}
	}
}
