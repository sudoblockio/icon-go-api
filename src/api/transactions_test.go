package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/metrics"
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
			description:  "Basic version",
			route:        "/api/v1/transactions",
			expectedCode: 200,
		},
		{
			description:  "Basic version",
			route:        "/api/v1/blocks",
			expectedCode: 200,
		},
	}

	metrics.Start()

	app := Start()

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		resp, err := app.Test(req, 1000)

		if err != nil {
			zap.Error(err)
		}

		assert.Equalf(t, nil, err, "app.Test(req)")
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
