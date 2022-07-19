package coincap

// Rate contains the exchange rate of a given asset in terms of USD as well as
// common identifiers for the asset in question and whether or not it is a fiat currency
type Rate struct {
	ID             string  `json:"id"`             // unique identifier for asset or fiat
	Symbol         string  `json:"symbol"`         // most common symbol used to identify asset or fiat
	CurrencySymbol string  `json:"currencySymbol"` // currency symbol if available
	RateUSD        float64 `json:"rateUsd,string"` // rate conversion to USD
	Type           string  `json:"type"`           // type of currency - fiat or crypto
}

// Rates returns currency rates standardized in USD.
func (c *Client) Rates() ([]Rate, Timestamp, error) {
	return request[[]Rate](c, "rates", nil)
}

// RateByID returns the USD rate for the given asset identifier.
func (c *Client) RateByID(id string) (*Rate, Timestamp, error) {
	return request[*Rate](c, "rates/"+id, nil)
}
