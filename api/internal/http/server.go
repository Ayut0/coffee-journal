package http

import (
	"net/http"

	"github.com/Ayut0/coffee-journal/api/builder"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer creates and configures an Echo instance with middleware and stub routes.
// It does not call Start() — that is the caller's responsibility.
func NewServer(d *builder.Dependency) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	stub := func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "not implemented"})
	}
	e.GET("/api/beans", stub)
	e.GET("/api/tastings", stub)
	e.GET("/api/search", stub)

	return e
}
