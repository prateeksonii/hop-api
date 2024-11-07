package main

import (
	"drop/db"
	"drop/lib"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}

	e := echo.New()
	lib.AddMiddlewares(e)

	db.InitDB()

	lib.InitRoutes(e)
	e.Logger.Panic(e.Start(":8000"))
}
