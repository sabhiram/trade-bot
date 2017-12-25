package types

// Balance keeps track of the available and total balance for a
// given currency.
type Balance struct {
	Currency  string
	Available float64
	Total     float64
}
