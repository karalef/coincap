package coincap

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"

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
func (c *Client) Trades(exchange string) (*Stream[*Trade], error) {
	e, _, err := c.ExchangeByID(exchange)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("exchange not found")
	}
	if !e.Socket {
		return nil, errors.New("exchange '" + exchange + "' does not support websockets")
	}
	const u = "wss://ws.coincap.io/trades/"
	return dial[*Trade](c.ws, u+exchange)
}

type price float64

// UnmarshalJSON is json.Unmarshaler implementation.
func (p *price) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(data[1:len(data)-1]), 64)
	*p = price(v)
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
func (c *Client) Prices(assets ...string) (*Stream[map[string]float64], error) {
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
	s, err := dial[map[string]price](c.ws, u+a)
	return (*Stream[map[string]float64])(unsafe.Pointer(s)), err
}

// Stream streams data from websocket conneсtion.
type Stream[T any] struct {
	ch   chan T
	stop chan struct{}
	conf chan struct{}
	err  error
}

// DataChannel returns data channel.
// It will be closed if there is an error or if the stream is closed.
func (s *Stream[T]) DataChannel() <-chan T {
	return s.ch
}

// Close closes stream.
func (s *Stream[T]) Close() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	<-s.conf
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) dial(conn *websocket.Conn) {
	defer func() {
		select {
		case <-s.stop:
			s.err = nil
		}
		conn.Close()
		close(s.ch)
		close(s.conf)
	}()
	for {
		_, r, err := conn.NextReader()
		if err != nil {
			s.err = err
			return
		}
		v, err := decodeJSON[T](r)
		if err != nil {
			s.err = err
			return
		}

		select {
		case <-s.stop:
			return
		case s.ch <- *v:
		}
	}
}

func dial[T any](ws *websocket.Dialer, u string) (*Stream[T], error) {
	conn, _, err := ws.Dial(u, nil)
	if err != nil {
		return nil, err
	}

	s := Stream[T]{
		ch:   make(chan T),
		stop: make(chan struct{}),
		conf: make(chan struct{}),
	}
	go s.dial(conn)

	return &s, nil
}
