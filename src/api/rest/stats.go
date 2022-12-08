package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sudoblockio/icon-go-api/config"
	"strconv"
)

func StatsAddHandlers(app *fiber.App) {
	prefix := config.Config.RestPrefix + "/stats"

	app.Get(prefix+"/", handlerGetStats)
	app.Get(prefix+"/circulating-supply", handlerGetCirculatingSupply)
	app.Get(prefix+"/market-cap", handlerGetMarketCap)
}

// Stats
// @Summary Get Stats
// @Description get json with a summary of stats
// @Tags Stats
// @BasePath /api/v1
// @Accept */*
// @Router /api/v1/stats [get]
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
func handlerGetStats(c *fiber.Ctx) error {
	UpdateCirculatingSupply()
	UpdateMarketCap()

	stats := map[string]interface{}{
		"circulating-supply": CirculatingSupply,
		"market-cap":         MarketCap,
	}
	body, _ := json.Marshal(stats)

	return c.SendString(string(body))
}

// Circulating Supply
// @Summary Get Stats
// @Description get circulating supply (total supply - burn wallet balance)
// @Tags Stats
// @BasePath /api/v1
// @Accept */*
// @Router /api/v1/stats/circulating-supply [get]
// @Success 200 {object} float64
// @Failure 422 {object} map[string]interface{}
func handlerGetCirculatingSupply(c *fiber.Ctx) error {
	UpdateCirculatingSupply()
	return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
}

// Market Cap
// @Summary Get Stats
// @Description get mkt cap (Coin Gecko Price * circulating supply)
// @Tags Stats
// @BasePath /api/v1
// @Accept */*
// @Router /api/v1/stats/market-cap [get]
// @Success 200 {object} float64
// @Failure 422 {object} map[string]interface{}
func handlerGetMarketCap(c *fiber.Ctx) error {
	UpdateMarketCap()
	return c.SendString(strconv.FormatFloat(MarketCap, 'f', -1, 64))
}
