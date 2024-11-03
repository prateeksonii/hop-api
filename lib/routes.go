package lib

import (
	"drop/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/ws", handlers.WsHandler)
	e.POST("/signup", handlers.SignUp)
	e.POST("/signin", handlers.SignIn)
	e.POST("/refresh", handlers.RefreshAuth)
}
