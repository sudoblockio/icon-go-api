package rest

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/crud"
	"github.com/sudoblockio/icon-go-api/redis"
)

// BlocksAddHandlers - add blocks endpoints to fiber router
func BlocksAddHandlers(app *fiber.App) {

	prefix := config.Config.RestPrefix + "/blocks"

	app.Get(prefix+"/", handlerGetBlocks)
	app.Get(prefix+"/:number", handlerGetBlockDetails)
	app.Get(prefix+"/timestamp/:timestamp", handlerGetBlockTimestampDetails)
}

// Parameters for handlerGetBlocks
type paramsGetBlocks struct {
	Limit       int    `query:"limit"`
	Skip        int    `query:"skip"`
	Number      uint32 `query:"number"`
	StartNumber uint32 `query:"start_number"`
	EndNumber   uint32 `query:"end_number"`
	Hash        string `query:"hash"`
	CreatedBy   string `query:"created_by"`
	Sort        string `query:"sort"`
}

// Blocks
// @Summary Get Blocks
// @Description get historical blocks
// @Tags Blocks
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param limit query int false "amount of records"
// @Param skip query int false "skip to a record"
// @Param number query int false "find by block number"
// @Param start_number query int false "range by start block number"
// @Param end_number query int false "range by end block number"
// @Param hash query string false "find by block hash"
// @Param created_by query string false "find by block creator"
// @Param sort query string false "desc or asc"
// @Router /api/v1/blocks [get]
// @Success 200 {object} []models.BlockList
// @Failure 422 {object} map[string]interface{}
func handlerGetBlocks(c *fiber.Ctx) error {
	params := &paramsGetBlocks{}
	if err := c.QueryParser(params); err != nil {
		zap.S().Warnf("Blocks Get Handler ERROR: %s", err.Error())

		c.Status(422)
		return c.SendString(`{"error": "could not parse query parameters"}`)
	}

	// Default params
	if params.Limit == 0 {
		params.Limit = 25
	}
	if params.Sort == "" {
		params.Sort = "desc"
	}

	// Check params
	if params.Limit < 1 || params.Limit > config.Config.MaxPageSize {
		c.Status(422)
		return c.SendString(`{"error": "invalid limit"}`)
	}
	if params.Skip < 0 || params.Skip > config.Config.MaxPageSkip {
		c.Status(422)
		return c.SendString(`{"error": "invalid skip"}`)
	}
	if params.EndNumber < params.StartNumber {
		c.Status(422)
		return c.SendString(`{"error": "end_number is less than start_number"}`)
	}
	if params.Sort != "desc" && params.Sort != "asc" {
		params.Sort = "desc"
	}

	blocks, err := crud.GetBlockCrud().SelectMany(
		params.Limit,
		params.Skip,
		params.Number,
		params.StartNumber,
		params.EndNumber,
		params.Hash,
		params.CreatedBy,
		params.Sort,
	)
	if err != nil {
		c.Status(500)
		zap.S().Warn(
			"Endpoint=handlerGetBlocks",
			" Error=Could not retrieve blocks: ", err.Error(),
		)
		return c.SendString(`{"error": "could not retrieve blocks"}`)
	}
	if len(*blocks) == 0 {
		// No Content
		c.Status(204)
	}

	// Set X-TOTAL-COUNT
	count, err := redis.GetRedisClient().GetCount(config.Config.RedisKeyPrefix + "block_count")
	if err != nil {
		count = 0
		zap.S().Warn(
			"Endpoint=handlerGetBlocks",
			" Error=Could not retrieve block count: ", err.Error(),
		)
	}
	c.Append("X-TOTAL-COUNT", strconv.FormatInt(count, 10))

	body, _ := json.Marshal(&blocks)
	return c.SendString(string(body))
}

// Block Details
// @Summary Get Block Details
// @Description get details of a block
// @Tags Blocks
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param number path int true "block number"
// @Router /api/v1/blocks/{number} [get]
// @Success 200 {object} models.Block
// @Failure 422 {object} map[string]interface{}
func handlerGetBlockDetails(c *fiber.Ctx) error {
	numberRaw := c.Params("number")

	if numberRaw == "" {
		c.Status(422)
		return c.SendString(`{"error": "number required"}`)
	}

	// Is number?
	number, err := strconv.ParseUint(numberRaw, 10, 32)
	if err != nil {
		c.Status(422)
		return c.SendString(`{"error": "invalid number"}`)
	}

	block, err := crud.GetBlockCrud().SelectOne(uint32(number))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return c.SendString(`{"error": "no block found"}`)
		}
		c.Status(500)
		return c.SendString(`{"error": "could not retrieve block"}`)
	}
	body, _ := json.Marshal(&block)
	return c.SendString(string(body))
}

// Block by Timestamp Details
// @Summary Get Block Details By Nearest Timestamp
// @Description get details of a block based on timestamp in millisecond epoch time
// @Tags Blocks
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param timestamp path int true "timestamp"
// @Router /api/v1/blocks/timestamp/{timestamp} [get]
// @Success 200 {object} models.Block
// @Failure 422 {object} map[string]interface{}
func handlerGetBlockTimestampDetails(c *fiber.Ctx) error {
	timestampRaw := c.Params("timestamp")

	if timestampRaw == "" {
		c.Status(422)
		return c.SendString(`{"error": "number required"}`)
	}

	// Is number?
	timestamp, err := strconv.ParseUint(timestampRaw, 10, 64)
	if err != nil {
		c.Status(422)
		return c.SendString(`{"error": "invalid number"}`)
	}

	block, err := crud.GetBlockCrud().SelectOneByTimestamp(uint64(timestamp))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return c.SendString(`{"error": "no block found"}`)
		}
		c.Status(500)
		return c.SendString(`{"error": "could not retrieve block"}`)
	}

	body, _ := json.Marshal(&block)
	return c.SendString(string(body))
}
