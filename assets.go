package coincap

import "net/url"

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

// Assets returns a list of all CoinCap assets.
func (c *Client) Assets() ([]Asset, Timestamp, error) {
	return c.AssetsSearch("", nil)
}

// AssetsSearch returns a list of CoinCap assets with params.
func (c *Client) AssetsSearch(search string, trim *TrimParams) ([]Asset, Timestamp, error) {
	q := make(url.Values)
	if search != "" {
		q.Set("search", search)
	}
	trim.setTo(&q)

	return request[[]Asset](c, "assets", q)
}

// AssetsSearchByIDs returns a list of CoinCap assets.
func (c *Client) AssetsSearchByIDs(ids []string) ([]Asset, Timestamp, error) {
	if ids == nil {
		return nil, 0, nil
	}
	return request[[]Asset](c, "assets", url.Values{"ids": ids})
}

// AssetByID returns an asset by its ID.
func (c *Client) AssetByID(id string) (*Asset, Timestamp, error) {
	return request[*Asset](c, "assets/"+id, nil)
}

// AssetHistory contains the USD price of an asset at a given timestamp.
type AssetHistory struct {
	PriceUSD float64   `json:"priceUsd,string"` // volume-weighted price in USD based on real-time market data
	Time     Timestamp `json:"time"`            // timestamp correlating to the given price
}

// AssetHistory returns USD price history of a given asset.
func (c *Client) AssetHistory(id string, interval *IntervalParams) ([]AssetHistory, Timestamp, error) {
	q := make(url.Values)
	var err = interval.setTo(&q, false)
	if err != nil {
		return nil, 0, err
	}
	return request[[]AssetHistory](c, "assets/"+id+"/history", q)
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

// AssetMarkets returns markets info of a given asset.
func (c *Client) AssetMarkets(id string, trim *TrimParams) ([]AssetMarket, Timestamp, error) {
	q := make(url.Values)
	trim.setTo(&q)
	return request[[]AssetMarket](c, "assets/"+id+"/markets", q)
}
