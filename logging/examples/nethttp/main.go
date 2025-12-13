package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/routerarchitects/ra-common-mods/logging"
	"github.com/routerarchitects/ra-common-mods/logging/middleware"
)

func main() {
	logging.InitService("example-nethttp")
	log := logging.Subsystem("http")
	app := fiber.New()

	// Logging middleware (adds request_id to ctxmeta + logs request)
	app.Use(middleware.RequestLogger(middleware.RequestLoggerOptions{
		// Default header is X-Request-Id; leave empty to use default
		EchoRequestID: true,
		LogEnd:        false,
		LogStart:      true,
	}))

	// Example route
	app.Get("/health", func(c *fiber.Ctx) error {
		log.Info(c.UserContext(), "health check", nil)
		return c.SendString("ok")
	})
	log.Info(context.TODO(), "starting server", logging.Fields{"addr": ":8080"})
	_ = app.Listen(":8086")
}
