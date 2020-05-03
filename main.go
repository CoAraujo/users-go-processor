package main

import (
	"context"
	"github.com/coaraujo/users-go-processor/infrastructure/config"
	"github.com/coaraujo/users-go-processor/infrastructure/queue"
	"github.com/coaraujo/users-go-processor/infrastructure/storage"
	"github.com/coaraujo/users-go-processor/processor"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := queue.GetInstance().NewConnection(); err != nil {
		log.Errorf("[Go-Processor] Fail to connect with ActiveMQ. Error: %s ", err)
	}
	defer queue.GetInstance().Disconnect()

	credential := options.Credential{
		Username:      config.MongodbUser,
		Password:      config.MongodbPassword,
		PasswordSet:   true,
		AuthSource:    config.MongodbDatabase,
		AuthMechanism: config.MongodbAuth,
	}
	if err := storage.GetInstance().Initialize(ctx, credential, "mongodb://"+config.MongodbHost+":"+config.MongodbPort,
		config.MongodbDatabase); err != nil {
		e.Logger.Fatal("[Go-Processor] Could not resolve Data access layer: ", err)
	}

	go queue.GetInstance().Listen(config.UserCreateTopic)
	go queue.GetInstance().Listen(config.UserRemovedTopic)
	go processor.GetInstance().Process()

	loadHealthcheck(e)
	setupServer(e)
}

func setupServer(e *echo.Echo) {
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "9090"
		}
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("[Go-Processor] Shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	queue.GetInstance().Disconnect()
}

func loadHealthcheck(e *echo.Echo) {
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "it's alive")
	})
}
