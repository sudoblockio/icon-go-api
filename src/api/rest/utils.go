package rest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strconv"
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

func respondWithCSV[T any](c *fiber.Ctx, data []T) error {
	var buf bytes.Buffer
	wr := csv.NewWriter(&buf)

	if len(data) == 0 {
		return c.Type("csv").Send(buf.Bytes())
	}

	// Use reflection to dynamically generate the CSV header
	val := reflect.ValueOf(data[0])
	var header []string
	for i := 0; i < val.Type().NumField(); i++ {
		header = append(header, val.Type().Field(i).Name)
	}
	if err := wr.Write(header); err != nil {
		return err
	}

	// Write data
	for _, item := range data {
		var record []string
		val := reflect.ValueOf(item)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			switch field.Kind() {
			case reflect.String:
				record = append(record, field.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				record = append(record, strconv.FormatInt(field.Int(), 10))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				record = append(record, strconv.FormatUint(field.Uint(), 10))
			default:
				record = append(record, fmt.Sprintf("%v", field.Interface()))
			}
		}
		if err := wr.Write(record); err != nil {
			return err
		}
	}
	wr.Flush()

	if err := wr.Error(); err != nil {
		return err
	}

	return c.Type("csv").Send(buf.Bytes())
}
