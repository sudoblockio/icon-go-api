package service

import (
	"errors"
	"fmt"
)

func IconNodeServiceGetTotalSupply() (float64, error) {

	// Request icon contract
	payload := fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "icx_getTotalSupply",
    "id": 1
	}`)

	body, err := JsonRpcRequestWithRetry(payload)
	if err != nil {
		return 0, err
	}

	// Extract result
	result, ok := body["result"].(string)
	if ok == false {
		return 0, errors.New("Cannot read result")
	}

	return StringHexToFloat64(result), nil
}

func IconNodeServiceGetBalance(publicKey string) (float64, error) {
	if publicKey == "hx0000000000000000000000000000000000000000" {
		return 0.0, nil
	} else if publicKey == "hx0000000000000000000000000000000000000001" {
		return 0.0, nil
	}

	payload := fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "icx_getBalance",
    "id": 1234,
    "params": {
        "address": "%s"
    }
	}`, publicKey)

	body, err := JsonRpcRequestWithRetry(payload)
	if err != nil {
		return 0.0, err
	}

	// Extract balance
	balance, ok := body["result"].(string)
	if ok == false {
		return 0.0, errors.New("Invalid response")
	}

	return StringHexToFloat64(balance), nil
}
