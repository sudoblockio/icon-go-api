package rest

import (
	"encoding/json"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/crud"
)

type AddressesQuery struct {
	Limit   int    `query:"limit"`
	Skip    int    `query:"skip"`
	Address string `query:"address"`
}

func AddressesAddHandlers(app *fiber.App) {

	prefix := config.Config.RestPrefix + "/addresses"

	app.Get(prefix+"/", handlerGetAddresses)
	app.Get(prefix+"/details/:address", handlerGetAddressDetails)
	app.Get(prefix+"/contracts", handlerGetContracts)
	app.Get(prefix+"/token-addresses/:address", handlerGetTokenAddresses)
}

// Addresses
// @Summary Get Addresses
// @Description get list of addresses
// @Tags Addresses
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param limit query int false "amount of records"
// @Param skip query int false "skip to a record"
// @Param is_contract query bool false "contract addresses only"
// @Param address query string false "find by address"
// @Router /api/v1/addresses [get]
// @Success 200 {object} []models.AddressList
// @Failure 422 {object} map[string]interface{}
func handlerGetAddresses(c *fiber.Ctx) error {
	params := new(AddressesQuery)
	if err := c.QueryParser(params); err != nil {
		zap.S().Warnf("Addresses Get Handler ERROR: %s", err.Error())

		c.Status(422)
		return c.SendString(`{"error": "could not parse query parameters"}`)
	}

	// Default Params
	if params.Limit <= 0 {
		params.Limit = 25
	}

	// Check Params
	if params.Limit < 1 || params.Limit > config.Config.MaxPageSize {
		c.Status(422)
		return c.SendString(`{"error": "limit must be greater than 0 and less than 101"}`)
	}
	if params.Skip < 0 || params.Skip > config.Config.MaxPageSkip {
		c.Status(422)
		return c.SendString(`{"error": "invalid skip"}`)
	}

	// Get Addresses
	addresses, err := crud.GetAddressCrud().SelectMany(
		params.Limit,
		params.Skip,
		params.Address,
	)
	if err != nil {
		zap.S().Warnf("Addresses CRUD ERROR: %s", err.Error())
		c.Status(500)
		return c.SendString(`{"error": "could not retrieve addresses"}`)
	}

	if len(*addresses) == 0 {
		// No Content
		c.Status(204)
	}

	// Set X-TOTAL-COUNT
	// TODO

	body, _ := json.Marshal(addresses)
	return c.SendString(string(body))
}

// Address Details
// @Summary Get Address Details
// @Description get details of an address
// @Tags Addresses
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param address path string true "find by address"
// @Router /api/v1/addresses/details/{address} [get]
// @Success 200 {object} models.Address
// @Failure 422 {object} map[string]interface{}
func handlerGetAddressDetails(c *fiber.Ctx) error {
	addressString := c.Params("address")
	if addressString == "" {
		c.Status(422)
		return c.SendString(`{"error": "address required"}`)
	}

	params := new(AddressesQuery)
	if err := c.QueryParser(params); err != nil {
		zap.S().Warnf("Addresses Get Handler ERROR: %s", err.Error())

		c.Status(422)
		return c.SendString(`{"error": "could not parse query parameters"}`)
	}

	// Get Addresses
	address, err := crud.GetAddressCrud().SelectOne(
		addressString,
	)
	if err != nil {
		c.Status(500)

		zap.S().Warnf("Addresses CRUD ERROR: %s", err.Error())
		return c.SendString(`{"error": "could not retrieve addresses"}`)
	}

	body, _ := json.Marshal(address)
	return c.SendString(string(body))
}

// Contract
// @Summary Get contracts
// @Description get list of contracts
// @Tags Addresses
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param limit query int false "amount of records"
// @Param skip query int false "skip to a record"
// @Router /api/v1/addresses/contracts [get]
// @Success 200 {object} []models.ContractList
// @Failure 422 {object} map[string]interface{}
func handlerGetContracts(c *fiber.Ctx) error {
	params := new(AddressesQuery)
	if err := c.QueryParser(params); err != nil {
		zap.S().Warnf("Addresses Get Handler ERROR: %s", err.Error())

		c.Status(422)
		return c.SendString(`{"error": "could not parse query parameters"}`)
	}

	// Default Params
	if params.Limit <= 0 {
		params.Limit = 25
	}

	// Check Params
	if params.Limit < 1 || params.Limit > config.Config.MaxPageSize {
		c.Status(422)
		return c.SendString(`{"error": "limit must be greater than 0 and less than 101"}`)
	}
	if params.Skip < 0 || params.Skip > config.Config.MaxPageSkip {
		c.Status(422)
		return c.SendString(`{"error": "invalid skip"}`)
	}

	// Get contracts
	contracts, err := crud.GetAddressCrud().SelectManyContracts(
		params.Limit,
		params.Skip,
	)
	if err != nil {
		zap.S().Warnf("Addresses CRUD ERROR: %s", err.Error())
		c.Status(500)
		return c.SendString(`{"error": "could not retrieve addresses"}`)
	}

	if len(*contracts) == 0 {
		// No Content
		c.Status(204)
	}

	// Set X-TOTAL-COUNT
	// TODO

	body, _ := json.Marshal(contracts)
	return c.SendString(string(body))
}

// Token Addresses
// @Summary Get Address Tokens
// @Description get list of token contracts by address
// @Tags Addresses
// @BasePath /api/v1
// @Accept */*
// @Produce json
// @Param address path string true "address"
// @Router /api/v1/addresses/token-addresses/{address} [get]
// @Success 200 {object} []string
// @Failure 422 {object} map[string]interface{}
func handlerGetTokenAddresses(c *fiber.Ctx) error {
	addressString := c.Params("address")

	// Get TokenAddresses
	tokenAddresss, err := crud.GetTokenAddressCrud().SelectManyByAddress(addressString)
	if err != nil {
		zap.S().Warnf("AddressTokens CRUD ERROR: %s", err.Error())
		c.Status(500)
		return c.SendString(`{"error": "could not retrieve addresses"}`)
	}

	if len(*tokenAddresss) == 0 {
		// No Content
		c.Status(204)
	}

	// Set X-TOTAL-COUNT
	// TODO

	// []models.TokenAddress -> []string
	var tokenContractAddresses []string
	for _, a := range *tokenAddresss {
		tokenContractAddresses = append(tokenContractAddresses, a.TokenContractAddress)
	}

	body, _ := json.Marshal(&tokenContractAddresses)
	return c.SendString(string(body))
}
