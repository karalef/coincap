package coincap

import (
	"errors"
	"net/url"
)

type candles struct{}

// Candle represets historic market performance for an asset over a given time span.
type Candle struct {
	Open   float64   `json:"open,string"`   // the price (quote) at which the first transaction was completed in a given time period
	High   float64   `json:"high,string"`   // the top price (quote) at which the base was traded during the time period
	Low    float64   `json:"low,string"`    // the bottom price (quote) at which the base was traded during the time period
	Close  float64   `json:"close,string"`  // the price (quote) at which the last transaction was completed in a given time period
	Volume float64   `json:"volume,string"` // the amount of base asset traded in the given time period
	Period Timestamp `json:"period"`        // timestamp for starting of that time period
}

// Candles returns all the market candle data for the provided exchange and parameters.
// The fields ExchangeID, BaseID, QuoteID, and Interval are required by the API.
func (candles) List(exchangeID, baseID, quoteID string, interval *IntervalParams, trim *TrimParams) ([]Candle, Timestamp, error) {
	// check required parameters.
	var err error
	if exchangeID == "" {
		err = errors.New("ExchangeID is required")
	} else if baseID == "" {
		err = errors.New("BaseID is required")
	} else if quoteID == "" {
		err = errors.New("QuoteID is required")
	}
	if err != nil {
		return nil, 0, err
	}

	var q = make(url.Values)
	q.Set("exchange", exchangeID)
	q.Set("baseId", baseID)
	q.Set("quoteId", quoteID)

	err = interval.set(&q, true)
	if err != nil {
		return nil, 0, err
	}
	trim.set(&q)

	var c []Candle
	return c, request(&c, "candles", q), nil
}
