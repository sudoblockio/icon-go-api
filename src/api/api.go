package api

import (
	"encoding/json"
	"strings"

	swagger "github.com/arsmn/fiber-swagger/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"

	_ "github.com/sudoblockio/icon-go-api/api/docs" // import for swagger docs
	"github.com/sudoblockio/icon-go-api/api/rest"
	"github.com/sudoblockio/icon-go-api/api/ws"
	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/global"
)

// @title Icon Go API
// @version 2.0
// @description The icon tracker API
func Start() {

	app := fiber.New()

	// Logging middleware
	app.Use(func(c *fiber.Ctx) error {
		zap.S().Info(c.Method(), " ", c.Path())

		// Go to next middleware:
		return c.Next()
	})

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:  config.Config.CORSAllowOrigins,
		AllowHeaders:  config.Config.CORSAllowHeaders,
		AllowMethods:  config.Config.CORSAllowMethods,
		ExposeHeaders: config.Config.CORSExposeHeaders,
	}))

	// Compression Middleware
	app.Use(compress.New(compress.Config{
		// refer to gofiber/fiber/blob/v1.14.6/middleware/compress.go#L17
		Level: compress.Level(config.Config.RestCompressLevel),
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Path(), "/docs/")
		},
	}))

	// Swagger docs
	app.Get(config.Config.RestPrefix+"/docs/*", swagger.Handler)

	// Add version handlers
	app.Get("/version", handlerVersion)
	app.Get("/metadata", handlerMetadata)

	// Add handlers
	rest.BlocksAddHandlers(app)
	rest.TransactionsAddHandlers(app)
	rest.LogsAddHandlers(app)
	rest.AddressesAddHandlers(app)
	ws.WebsocketsAddHandlers(app)

	go app.Listen(":" + config.Config.APIPort)
}

// Version
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Version
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /version [get]
func handlerVersion(c *fiber.Ctx) error {
	message := map[string]string{
		"version": global.Version,
	}

	jsonMessage, _ := json.Marshal(message)

	return c.SendString(string(jsonMessage))
}

// Metadata
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Version
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metadata [get]
func handlerMetadata(c *fiber.Ctx) error {
	message := map[string]string{
		"version":     global.Version,
		"name":        config.Config.Name,
		"description": "a go api template",
	}

	jsonMessage, _ := json.Marshal(message)

	return c.SendString(string(jsonMessage))
}
