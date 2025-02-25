package middleware

import (
	"crop_connect/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := helper.GetPayloadFromToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, helper.BaseResponse{
					Status:  http.StatusUnauthorized,
					Message: "token tidak valid",
					Data:    nil,
				})
			}

			return next(c)
		}
	}
}

func CheckOneRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := helper.GetPayloadFromToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, helper.BaseResponse{
					Status:  http.StatusUnauthorized,
					Message: "token tidak valid",
					Data:    nil,
				})
			}

			if token.Role == role {
				return next(c)
			}

			return c.JSON(http.StatusForbidden, helper.BaseResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden",
				Data:    nil,
			})
		}
	}
}

func CheckManyRole(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := helper.GetPayloadFromToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, helper.BaseResponse{
					Status:  http.StatusUnauthorized,
					Message: "token tidak valid",
					Data:    nil,
				})
			}

			for _, role := range roles {
				if token.Role == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, helper.BaseResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden",
				Data:    nil,
			})
		}
	}
}
