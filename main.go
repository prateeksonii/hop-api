package main

import (
	"drop/db"
	"drop/lib"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	lib.AddMiddlewares(e)

	db.InitDB()

	lib.InitRoutes(e)
	e.Logger.Panic(e.Start(":8000"))
}
