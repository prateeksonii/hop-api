package lib

import (
	"drop/handlers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/ws", handlers.WsHandler)

	g := e.Group("/auth")
	g.POST("/signup", handlers.SignUp)
	g.POST("/signin", handlers.SignIn)
	g.POST("/refresh", handlers.RefreshAuth)

	g = e.Group("/api")
	g.Use(echojwt.JWT([]byte("restinpeace")))
	g.GET("/me", handlers.GetAuthenticatedUser)
}
