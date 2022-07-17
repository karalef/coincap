package coincap

import (
	"context"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// Trade type.
type Trade struct {
	Exchange  string  `json:"exchange"`
	Base      string  `json:"base"`
	Quote     string  `json:"quote"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
	PriceUSD  float64 `json:"priceUds"`
}

// Trades streams trades from other cryptocurrency exchange websockets.
// Users must select a specific exchange. In the /exchanges endpoint users
// can determine if an exchange has a socket available by noting
// response 'socket':true/false.
// The trades websocket is the only way to receive individual
// trade data through CoinCap.
func (c *Client) Trades(ctx context.Context, exchange string, ch chan<- *Trade) error {
	const u = "wss://ws.coincap.io/trades/"
	return dial(ctx, c.ws, u+exchange, ch)
}

// Price implements Unmarshaler interface for float64.
type Price float64

// UnmarshalJSON is Unmarshaler implementation.
func (p *Price) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(data), 64)
	*p = Price(v)
	return err
}

// Prices is the most accurate source of real-time changes to the global price
// of an asset. Each time the system receives data that moves the global price
// in one direction or another, this change is immediately published through
// the websocket.
// These prices correspond with the values shown in /assets - a value that may
// change several times per second based on market activity.
// Emtpy 'assets' means prices for all assets.
func (c *Client) Prices(ctx context.Context, ch chan<- *map[string]Price, assets ...string) error {
	a := "ALL"
	if len(assets) > 0 {
		a = strings.Join(assets, ",")
	}
	const u = "wss://ws.coincap.io/prices?assets="
	return dial(ctx, c.ws, u+a, ch)
}

func dial[T any](ctx context.Context, ws *websocket.Dialer, u string, ch chan<- *T) error {
	conn, _, err := ws.DialContext(ctx, u, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		var v T
		err = conn.ReadJSON(&v)
		if err != nil {
			break
		}
		ch <- &v
	}
	return err
}
