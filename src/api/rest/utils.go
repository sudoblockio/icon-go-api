package rest

import (
	"github.com/sudoblockio/icon-go-api/service"
)

// https://stackoverflow.com/a/28818489/12642712
func newTrue() *bool {
	b := true
	return &b
}

var addressSortParams = []string{"name", "balance", "transaction_count", "transaction_internal_count", "token_transfer_count"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getCirculatingSupply() (float64, error) {
	totalSupply, err := service.IconNodeServiceGetTotalSupply()
	if err != nil {
		return 0, err
	}

	burnBalance, err := service.IconNodeServiceGetBalance("hx1000000000000000000000000000000000000000")
	if err != nil {
		return 0, err
	}
	circulatingSupply := totalSupply - burnBalance
	return circulatingSupply, err
}
