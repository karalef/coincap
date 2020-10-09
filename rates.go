package coincap

type rates struct{}

// Rate contains the exchange rate of a given asset in terms of USD as well as
// common identifiers for the asset in question and whether or not it is a fiat currency
type Rate struct {
	ID             string  `json:"id"`             // unique identifier for asset or fiat
	Symbol         string  `json:"symbol"`         // most common symbol used to identify asset or fiat
	CurrencySymbol string  `json:"currencySymbol"` // currency symbol if available
	RateUSD        float64 `json:"rateUsd,string"` // rate conversion to USD
	Type           string  `json:"type"`           // type of currency - fiat or crypto
}

// List returns currency rates standardized in USD.
func (rates) List() ([]Rate, Timestamp) {
	var r []Rate
	return r, request(&r, "rates", nil)
}

// ByID returns the USD rate for the given asset identifier.
func (rates) ByID(id string) (*Rate, Timestamp) {
	var r Rate
	return &r, request(&r, "rates/"+id, nil)
}
