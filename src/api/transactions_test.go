package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/sudoblockio/icon-go-api/config"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

func TestRestTransactions(t *testing.T) {
	config.ReadEnvironment()

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "txs",
			route:        "/api/v1/transactions",
			expectedCode: 200,
		},
		{
			description:  "token-transfers",
			route:        "/api/v1/transactions/token-transfers",
			expectedCode: 200,
		},
	}

	app := Start()

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		resp, err := app.Test(req, 10000)

		if err != nil {
			zap.Error(err)
		}

		assert.Equalf(t, nil, err, "app.Test(req)")
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
