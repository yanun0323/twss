package app

import (
	"context"
	"errors"
	"net/http"
	"stocker/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yanun0323/pkg/logs"
)

var (
	_RateLimit *echo.MiddlewareFunc
)

func DefaultMiddleware(ctx context.Context) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		rateLimit(ctx),
	}
}

func rateLimit(ctx context.Context) echo.MiddlewareFunc {
	if _RateLimit != nil {
		return *_RateLimit
	}
	l := logs.Get(ctx)
	r := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store:   middleware.NewRateLimiterMemoryStore(1),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			if len(id) == 0 {
				return id, errors.New("empty ip address")
			}
			return id, nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			l.Warnf("[%s] rate limit exceeded, %+v", c.RealIP(), err)
			return c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid request", err))
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			l.Warnf("[%s] rate limit exceeded", identifier)
			return c.JSON(
				http.StatusTooManyRequests,
				util.NewErrorResponse("rate limit exceeded, you should not send more than 1 request per second ", err),
			)
		}})
	_RateLimit = &r
	return *_RateLimit
}
