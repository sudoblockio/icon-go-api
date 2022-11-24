package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringHexToInt64(t *testing.T) {
	inputString := "0x311686fe637dc7b0622d7e6"
	output := StringHexToFloat64(inputString)
	assert.Equal(t, float64(949499958.6892647), output, "Test success")
}
