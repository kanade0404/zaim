package main

import (
	"github.com/labstack/echo/v4"
	"zaim/handlers/zaim"
	"zaim/middlewares"
)

func main() {
	e := echo.New()
	initialize(e)
	e.Logger.Fatal(e.Start(":8888"))
}
func initialize(e *echo.Echo) {
	e.Use(middlewares.Context)
	auth := e.Group("/auth")
	auth.GET("/", zaim.Authorization)
	auth.GET("/callback", zaim.CallbackOAuthToken)
	e.POST("/transaction", zaim.Register)
	e.GET("/categories", zaim.ListActiveCategory)
	e.GET("/genres", zaim.ListActiveGenre)
	e.GET("/accounts", zaim.ListActiveAccount)
}
