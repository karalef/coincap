package coincap

import "net/url"

type assets struct{}

// Asset contains CoinCap asset data from exchanges.
type Asset struct {
	ID                string  `json:"id"`                       // unique identifier for asset
	Rank              int     `json:"rank,string"`              // rank in terms of the asset's market cap
	Symbol            string  `json:"symbol"`                   // most common symbol used to identify this asset on an exchange
	Name              string  `json:"name"`                     // proper name for asset
	Supply            float64 `json:"supply,string"`            // available supply for trading
	MaxSupply         float64 `json:"maxSupply,string"`         // total quantity of asset issued
	MarketCapUsd      float64 `json:"marketCapUsd,string"`      // supply x price
	VolumeUsd24Hr     float64 `json:"volumeUsd24Hr,string"`     // quantity of trading volume in USD over last 24 hours
	PriceUsd          float64 `json:"priceUsd,string"`          // volume-weighted price based on real-time market data, translated to USD
	ChangePercent24Hr float64 `json:"changePercent24Hr,string"` // the direction and value change in the last 24 hours
	Vwap24Hr          float64 `json:"vwap24Hr,string"`          // Volume Weighted Average Price in the last 24 hours
}

// List returns a list of all CoinCap Assets.
func (assets) List() ([]Asset, Timestamp) {
	return Assets.Get("", nil, nil)
}

// Get returns a list of CoinCap Assets with params.
func (assets) Get(search string, ids []string, params *TrimParams) ([]Asset, Timestamp) {
	var q = make(url.Values)
	if ids != nil {
		for _, id := range ids {
			q.Add("ids", id)
		}
	} else {
		if search != "" {
			q.Set("search", search)
		}
		params.set(&q)
	}

	var list []Asset
	return list, request(&list, "assets", q)
}

// GetByID returns an asset by its ID.
func (assets) ByID(id string) (*Asset, Timestamp) {
	var a Asset
	return &a, request(&a, "assets/"+id, nil)
}

// AssetHistory contains the USD price of an asset at a given timestamp.
type AssetHistory struct {
	PriceUSD float64   `json:"priceUsd,string"` // volume-weighted price in USD based on real-time market data
	Time     Timestamp `json:"time"`            // timestamp correlating to the given price
}

// History returns USD price history of a given asset.
func (assets) History(id string, params *IntervalParams) ([]AssetHistory, Timestamp, error) {
	var q = make(url.Values)
	if err := params.set(&q, false); err != nil {
		return nil, 0, err
	}

	var history []AssetHistory
	return history, request(&history, "assets/"+id+"/history", q), nil
}

// AssetMarket contains the markets info of an asset.
type AssetMarket struct {
	ExchangeID    string  `json:"exchangeId"`           // unique identifier for exchange
	BaseID        string  `json:"baseId"`               // unique identifier for this asset, base is asset purchased
	QuoteID       string  `json:"quoteId"`              // unique identifier for this asset, quote is asset used to purchase based
	BaseSymbol    string  `json:"baseSymbol"`           // most common symbol used to identify asset, base is asset purchased
	QuoteSymbol   string  `json:"quoteSymbol"`          // most common symbol used to identify asset, quote is asset used to purchase base
	VolumeUsd24Hr float64 `json:"volumeUsd24Hr,string"` // volume transacted on this market in last 24 hours
	PriceUsd      float64 `json:"priceUsd,string"`      // the amount of quote asset traded for one unit of base asset
	VolumePercent float64 `json:"volumePercent,string"` // percent of quote asset volume
}

// Markets returns markets info of a given asset.
func (assets) Markets(id string, params *TrimParams) ([]AssetMarket, Timestamp) {
	var q = make(url.Values)
	params.set(&q)

	var m []AssetMarket
	return m, request(&m, "assets/"+id+"/markets", q)
}
