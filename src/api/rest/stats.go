package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sudoblockio/icon-go-api/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StatsAddHandlers(app *fiber.App) {
	prefix := config.Config.RestPrefix + "/stats"

	app.Get(prefix+"/circulating-supply", handlerGetCirculatingSupply)
	app.Get(prefix+"/market-cap", handlerGetMarketCap)
}

var CirculatingSupply float64
var LastUpdatedTimeCirculatingSupply time.Time

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
	timeDiff := time.Now().Sub(LastUpdatedTimeCirculatingSupply)
	if timeDiff > 1*time.Hour {
		circulatingSupply, err := getCirculatingSupply()
		if err != nil {
			return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
		}
		CirculatingSupply = circulatingSupply
		LastUpdatedTimeCirculatingSupply = time.Now()
	}
	return c.SendString(strconv.FormatFloat(CirculatingSupply, 'f', -1, 64))
}

var MarketCap float64
var LastUpdatedTimeMarketCap time.Time

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
	timeDiff := time.Now().Sub(LastUpdatedTimeMarketCap)
	if timeDiff > 1*time.Hour {
		resp, err := http.Get("https://api.coingecko.com/api/v3/coins/icon")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.SendString(strconv.FormatFloat(MarketCap, 'f', -1, 64))
		}

		response := make(map[string]interface{})
		err = json.Unmarshal(body, &response)
		if err != nil {
			return c.SendString(strconv.FormatFloat(MarketCap, 'f', -1, 64))
		}
		usdPrice, ok := response["market_data"].(map[string]interface{})["current_price"].(map[string]interface{})["usd"].(float64)
		if !ok {
			return nil
		}
		if CirculatingSupply == 0.0 {
			circulatingSupply, err := getCirculatingSupply()
			if err != nil {
				return c.SendString(strconv.FormatFloat(MarketCap, 'f', -1, 64))
			}
			CirculatingSupply = circulatingSupply
			LastUpdatedTimeCirculatingSupply = time.Now()
		}

		MarketCap = CirculatingSupply * usdPrice
		LastUpdatedTimeMarketCap = time.Now()
	}
	return c.SendString(strconv.FormatFloat(MarketCap, 'f', -1, 64))
}
