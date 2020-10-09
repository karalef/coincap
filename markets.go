package coincap

import "net/url"

type markets struct{}

// MarketsRequest contains the paramters you can use to provide a request for market data.
type MarketsRequest struct {
	ExchangeID  string // search by unique exchange ID
	BaseSymbol  string // return all results with this base symbol
	BaseID      string // return all results with this base id
	QuoteSymbol string // return all results with this quote symbol
	QuoteID     string // return all results with this quote ID
	AssetSymbol string // return all results with this asset symbol
	AssetID     string // return all results with this asset ID
}

// Market contains the market data response from the api.
type Market struct {
	ExchangeID            string    `json:"exchangeId"`                   // unique identifier for exchange
	Rank                  int       `json:"rank,string"`                  // rank in terms of volume transacted in this market
	BaseSymbol            string    `json:"baseSymbol"`                   // most common symbol used to identify this asset
	BaseID                string    `json:"baseId"`                       // unique identifier for this asset. base is the asset purchased
	QuoteSymbol           string    `json:"quoteSymbol"`                  // most common symbol used to identify this asset
	QuoteID               string    `json:"quoteId"`                      // unique identifier for thisasset. quote is the asset used to purchase base
	PriceQuote            float64   `json:"priceQuote,string"`            // amount of quote asset traded for 1 unit of base asset
	PriceUsd              float64   `json:"priceUsd,string"`              // quote price translated to USD
	VolumeUsd24Hr         float64   `json:"volumeUsd24Hr,string"`         // volume transacted in this market in the last 24 hours
	PercentExchangeVolume float64   `json:"percentExchangeVolume,string"` // amount of daily volume this market transacts compared to others on this exchange
	TradesCount24Hr       int       `json:"tradesCount24Hr,string"`       // number of trades on this market in the last 24 hours
	Updated               Timestamp `json:"updated"`                      // last time information was received from this market
}

// Markets requests market data for all markets matching the criteria set in the MarketRequest params.
func (markets) List(marketsParams MarketsRequest, params *TrimParams) ([]Market, Timestamp) {
	var q = make(url.Values)
	params.set(&q)
	if marketsParams.ExchangeID != "" {
		q.Set("exchange", marketsParams.ExchangeID)
	}
	if marketsParams.BaseSymbol != "" {
		q.Set("baseSymbol", marketsParams.BaseSymbol)
	}
	if marketsParams.BaseID != "" {
		q.Set("baseId", marketsParams.BaseID)
	}
	if marketsParams.QuoteSymbol != "" {
		q.Set("quoteSymbol", marketsParams.QuoteSymbol)
	}
	if marketsParams.QuoteID != "" {
		q.Set("quoteId", marketsParams.QuoteID)
	}
	if marketsParams.AssetSymbol != "" {
		q.Set("assetSymbol", marketsParams.AssetSymbol)
	}
	if marketsParams.AssetID != "" {
		q.Set("assetId", marketsParams.AssetID)
	}

	var m []Market
	return m, request(&m, "markets", q)
}
