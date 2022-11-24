package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/service"
	"strconv"
	"time"
)

func StatsAddHandlers(app *fiber.App) {
	prefix := config.Config.RestPrefix + "/stats"

	app.Get(prefix+"/circulating-supply", handlerGetCirculatingSupply)
}

var CirculatingSupply float64
var LastUpdatedTimeCirculatingSupply time.Time

// Circulating Supply
// @Summary Get Stats
// @Description get stats for ICON
// @Tags Stats
// @BasePath /api/v1
// @Accept */*
// @Router /api/v1/stats/circulating-supply [get]
// @Success 200 {object} []models.TransactionList
// @Failure 422 {object} map[string]interface{}
func handlerGetCirculatingSupply(c *fiber.Ctx) error {
	timeDiff := time.Now().Sub(LastUpdatedTimeCirculatingSupply)
	if timeDiff > 1*time.Hour {
		totalSupply, err := service.IconNodeServiceGetTotalSupply()
		if err != nil {
			return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
		}

		burnBalance, err := service.IconNodeServiceGetBalance("hx1000000000000000000000000000000000000000")
		if err != nil {
			return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
		}
		CirculatingSupply = totalSupply - burnBalance
		LastUpdatedTimeCirculatingSupply = time.Now()
	}
	return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
}
