package middleware

import (
	"errors"
	"main/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func TokenAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParams().Get("username")
		password := c.QueryParams().Get("password")
		if Valid(username, password) {
			return next(c)
		}
		return c.JSON(http.StatusBadRequest, response.Error(errors.New("wrong username and password")))
	}
}

func Valid(username, password string) bool {
	u := viper.GetString("auth.username")
	p := viper.GetString("auth.password")
	check := len(u) > 0 && len(p) > 0
	return check && username == u && password == p
}
