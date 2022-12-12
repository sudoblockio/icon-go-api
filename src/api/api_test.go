package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/sudoblockio/icon-go-api/config"
	"net/http/httptest"
	"testing"
)

func TestRestApi(t *testing.T) {
	config.ReadEnvironment()

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Docs",
			route:        "/api/v1/docs/index.html",
			expectedCode: 200,
		},
		{
			description:  "Version",
			route:        "/version",
			expectedCode: 200,
		},
		{
			description:  "get HTTP status 404, when route is not exists",
			route:        "/not-found",
			expectedCode: 404,
		},
		{
			description:  "addresses",
			route:        "/api/v1/addresses",
			expectedCode: 200,
		},
		{
			description:  "addresses contracts",
			route:        "/api/v1/addresses/contracts",
			expectedCode: 200,
		},
		{
			description:  "txs",
			route:        "/api/v1/transactions",
			expectedCode: 200,
		},
		//{
		//	description:  "txs internal",
		//	route:        "/api/v1/transactions/internal",
		//	expectedCode: 200,
		//},
		{
			description:  "token-transfers",
			route:        "/api/v1/transactions/token-transfers",
			expectedCode: 200,
		},
		{
			description:  "logs",
			route:        "/api/v1/logs",
			expectedCode: 200,
		},
	}

	app := Start()

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		resp, err := app.Test(req, 10000)
		assert.Equalf(t, nil, err, "app.Test(req)")
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
