package lib

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddMiddlewares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowHeaders:     middleware.DefaultCORSConfig.AllowHeaders,
		AllowMethods:     middleware.DefaultCORSConfig.AllowMethods,
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))
}
