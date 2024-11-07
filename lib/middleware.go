package lib

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func AddMiddlewares(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New(validator.WithRequiredStructEnabled())}

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
