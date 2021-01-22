package coincap

// Exchange contains information about a cryptocurrency exchange.
type Exchange struct {
	ID                 string    `json:"exchangeId"`                // unique identifier for exchange
	Name               string    `json:"name"`                      // proper name of exchange
	Rank               int       `json:"rank,string"`               // rank in terms of total volume compared to other exchanges
	PercentTotalVolume float64   `json:"percentTotalVolume,string"` // perecent of total daily volume in relation to all exchanges
	VolumeUSD          float64   `json:"volumeUSD,string"`          // daily volume represented in USD
	TradingPairs       int       `json:"tradingPairs,string"`       // number of trading pairs offered by the exchange
	Socket             bool      `json:"socket"`                    // Whether or not a trade socket is available on this exchange
	URL                string    `json:"exchangeUrl"`               // url of exchange
	Updated            Timestamp `json:"updated"`                   // Time since information was last updated
}

// Exchanges returns information about all exchanges currently tracked by CoinCap.
func (c *Client) List() ([]Exchange, Timestamp, error) {
	var e []Exchange
	ts, err := c.request(&e, "exchanges", nil)
	return e, ts, err
}

// ByID returns exchange data for an exchange with the given unique ID.
func (c *Client) ExchangeByID(id string) (*Exchange, Timestamp, error) {
	var e Exchange
	ts, err := c.request(&e, "exchanges/"+id, nil)
	return &e, ts, err
}
