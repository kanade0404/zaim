package main

import (
	"github.com/alexlast/bunzap"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/handlers/zaim"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"os"
)

func main() {
	initialize()
}

func initialize() {
	e := echo.New()
	e.Use(middleware.Logger())
	var (
		zapLogger *zap.Logger
		err       error
	)
	if os.Getenv("ENV") == "local" {
		zapLogger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalf("failed to create zap logger: %v", err)
		}
	} else {
		zapLogger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("failed to create zap logger: %v", err)
		}
	}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			zapLogger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
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
	api.POST("/category", zaim.UpdateCategory)
	api.POST("/genre", zaim.UpdateGenres)
	api.POST("/account", zaim.UpdateAccount)
	e.Logger.Fatal(e.Start(":8888"))
}
