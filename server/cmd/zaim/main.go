package main

import (
	"github.com/alexlast/bunzap"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/handlers/zaim"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"os"
)

func main() {
	initialize()
}

// Context echo.Context をラップする構造体を定義する
type Context struct {
	echo.Context
}

//// BindValidate Bind と Validate を合わせたメソッド
//func (c *Context) BindValidate(i interface{}) error {
//	if err := c.Bind(i); err != nil {
//		return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
//	}
//	if err := c.Validate(i); err != nil {
//		return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
//	}
//	return nil
//}
//
//type callFunc func(c *Context) error
//
//func c(h callFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		return h(c.(*Context))
//	}
//}

func initialize() {
	e := echo.New()
	e.Use(middleware.Logger())
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{c})
		}
	})
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	db.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger: zapLogger,
	}))
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	api := e.Group("/api")
	api.POST("/transaction/b43", zaim.B43Register)
	//api.GET("/category", zaim.ListActiveCategory)
	api.POST("/category", zaim.UpdateCategory)
	//api.GET("/genre", zaim.ListActiveGenre)
	api.POST("/genre", zaim.UpdateGenres)
	//api.GET("/account", zaim.ListActiveAccount)
	api.POST("/account", zaim.UpdateAccount)
	e.Logger.Fatal(e.Start(":8888"))
}
