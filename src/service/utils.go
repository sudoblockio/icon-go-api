package service

import (
	"go.uber.org/zap"
	"math/big"
)

func StringHexToFloat64(hex string) float64 {
	valueDecimal := float64(0)

	var negative bool
	if hex[:1] == "-" {
		hex = hex[1:]
		negative = true
	}

	valueBigInt, success := new(big.Int).SetString(hex[2:], 16)
	if success == false {
		zap.S().Warn("Set String Error: hex=", hex)
		return 0
	}

	baseBigFloatString := "1"
	for i := 0; i < 18; i++ {
		baseBigFloatString += "0"
	}
	baseBigFloat, success := new(big.Float).SetString(baseBigFloatString) // 10^(base)
	if success == false {
		zap.S().Warn("Set String Error: base=", 18)
		return 0
	}

	valueBigFloat := new(big.Float).SetInt(valueBigInt)
	valueBigFloat = valueBigFloat.Quo(valueBigFloat, baseBigFloat)

	valueDecimal, _ = valueBigFloat.Float64()

	if negative {
		valueDecimal = -1 * valueDecimal
	}

	return valueDecimal
}

//func StringHexToInt64(i string) int64 {
//	o := new(big.Int)
//
//	o.SetString(i, 0)
//	fmt.Println(o)
//
//	num := o.Int64()
//	fmt.Println(n)
//	return o.Int64()
//}
