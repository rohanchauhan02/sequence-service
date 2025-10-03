package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
)

func MiddlewareRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	}
}
