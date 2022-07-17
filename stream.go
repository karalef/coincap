package coincap

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"

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
func (c *Client) Trades(exchange string) (*Stream[Trade], error) {
	e, _, err := c.ExchangeByID(exchange)
	if err != nil {
		return nil, err
	}
	if !e.Socket {
		return nil, errors.New("exchange '" + exchange + "' does not support websockets")
	}
	const u = "wss://ws.coincap.io/trades/"
	return dial[Trade](c.ws, u+exchange)
}

// Price implements Unmarshaler interface for float64.
type Price float64

// UnmarshalJSON is Unmarshaler implementation.
func (p *Price) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(data[1:len(data)-1]), 64)
	*p = Price(v)
	return err
}

// Prices is the most accurate source of real-time changes to the global price
// of an asset. Each time the system receives data that moves the global price
// in one direction or another, this change is immediately published through
// the websocket.
// These prices correspond with the values shown in /assets - a value that may
// change several times per second based on market activity.
//
// Emtpy 'assets' means prices for all assets.
func (c *Client) Prices(assets ...string) (*Stream[map[string]Price], error) {
	a := "ALL"
	if len(assets) > 0 {
		t, _, err := c.AssetsSearchByIDs(assets)
		if err != nil {
			return nil, err
		}
		if len(t) != len(assets) {
			return nil, errors.New("incorrect assets ids")
		}
		a = strings.Join(assets, ",")
	}
	const u = "wss://ws.coincap.io/prices?assets="
	return dial[map[string]Price](c.ws, u+a)
}

// Stream streams data from websocket conneсtion.
type Stream[T any] struct {
	conn *websocket.Conn
	ch   chan *T
	stop chan struct{}
	mut  sync.Mutex
	err  error
}

// DataChannel returns data channel.
// It will be closed if there is an error or if the stream is closed.
func (s *Stream[T]) DataChannel() <-chan *T {
	return s.ch
}

// Close closes stream.
func (s *Stream[T]) Close() {
	s.mut.Lock()
	s.conn.Close()
	close(s.stop)
	if s.err == nil {
		close(s.ch)
	}
	s.mut.Unlock()
}

func (s *Stream[T]) Err() error {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.err
}

func (s *Stream[T]) dial() {
	var err error
	for {
		var r io.Reader
		_, r, err = s.conn.NextReader()
		if err != nil {
			break
		}
		var v *T
		var b []byte
		v, b, err = decodeJSON[T](r)
		if err != nil {
			err = errors.New(string(b))
			break
		}

		select {
		case <-s.stop:
			return
		case s.ch <- v:
		}
	}
	s.mut.Lock()
	select {
	case <-s.stop:
	default:
		s.err = err
		close(s.ch)
	}
	s.mut.Unlock()
	return
}

func dial[T any](ws *websocket.Dialer, u string) (*Stream[T], error) {
	conn, _, err := ws.Dial(u, nil)
	if err != nil {
		return nil, err
	}

	s := Stream[T]{
		conn: conn,
		ch:   make(chan *T),
		stop: make(chan struct{}),
	}
	go s.dial()

	return &s, nil
}
