package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/sudoblockio/icon-go-api/config"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

func TestRestStats(t *testing.T) {
	config.ReadEnvironment()

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "supply",
			route:        "/api/v1/stats/circulating-supply",
			expectedCode: 200,
		},
		{
			description:  "mkt cap",
			route:        "/api/v1/stats/market-cap",
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
