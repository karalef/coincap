package coincap

import (
	"testing"
)

func TestTradesInvalidExchange(t *testing.T) {
	c := DefaultClient

	_, err := c.Trades("zxcvzxvzv")
	if err == nil {
		t.Fail()
	}
}

func TestTrades(t *testing.T) {
	c := DefaultClient

	// make sure that binance supports websockets
	e, _, err := c.ExchangeByID("binance")
	if err != nil {
		t.Fatal(err)
	}
	if e == nil {
		t.Fatal("binance exchange not found")
	}
	if !e.Socket {
		t.Fatal("binance does not support websockets")
	}

	s, err := c.Trades("binance")
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	ch := s.DataChannel()
	for {
		_, ok := <-ch
		if count == 10 {
			if ok {
				t.Fatal("stream did not close")
			}
			break
		}
		if !ok {
			t.Fatal(s.Err())
		}

		count++
		if count == 10 {
			s.Close()
		}
	}
}

func TestPricesInvalidAsset(t *testing.T) {
	c := DefaultClient

	_, err := c.Prices("asdasdasasdasd")
	if err == nil {
		t.Fail()
	}
}

func TestPrices(t *testing.T) {
	c := DefaultClient

	s, err := c.Prices("bitcoin")
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	ch := s.DataChannel()
	for {
		d, ok := <-ch
		if count == 3 {
			if ok {
				t.Fatal("stream did not close")
			}
			break
		}
		if !ok {
			t.Fatal(s.Err())
		}

		_, ok = (*d)["bitcoin"]
		if !ok {
			t.Fatal("result has no expected field")
		}

		count++
		if count == 3 {
			s.Close()
		}
	}
}
