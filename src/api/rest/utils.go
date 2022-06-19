package rest

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
