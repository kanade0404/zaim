package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"zaim/handlers/zaim"
	"zaim/middlewares"
)

func main() {
	initialize()
}
func initialize() {
	e := echo.New()
	e.Use(middlewares.Context)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)
	auth := e.Group("/auth")
	auth.GET("/", zaim.Authorization)
	auth.GET("/callback", zaim.CallbackOAuthToken)
	e.POST("/transaction", zaim.Register)
	e.GET("/categories", zaim.ListActiveCategory)
	e.GET("/genres", zaim.ListActiveGenre)
	e.GET("/accounts", zaim.ListActiveAccount)
	e.Logger.Fatal(e.Start(":8888"))
}
