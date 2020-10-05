package coincap

type exchanges struct{}

// Exchange contains information about a cryptocurrency exchange.
type Exchange struct {
	ID                 string    `json:"exchangeId"`         // unique identifier for exchange
	Name               string    `json:"name"`               // proper name of exchange
	Rank               Int       `json:"rank"`               // rank in terms of total volume compared to other exchanges
	PercentTotalVolume Float     `json:"percentTotalVolume"` // perecent of total daily volume in relation to all exchanges
	VolumeUSD          Float     `json:"volumeUSD"`          // daily volume represented in USD
	TradingPairs       Int       `json:"tradingPairs"`       // number of trading pairs offered by the exchange
	Socket             bool      `json:"socket"`             // Whether or not a trade socket is available on this exchange
	URL                string    `json:"exchangeUrl"`        // url of exchange
	Updated            Timestamp `json:"updated"`            // Time since information was last updated
}

// List returns information about all exchanges currently tracked by CoinCap.
func (exchanges) List() ([]Exchange, Timestamp) {
	var e []Exchange
	return e, request(&e, "exchanges", nil)
}

// ByID returns exchange data for an exchange with the given unique ID.
func (exchanges) ByID(id string) (*Exchange, Timestamp) {
	var e Exchange
	return &e, request(&e, "exchanges/"+id, nil)
}
