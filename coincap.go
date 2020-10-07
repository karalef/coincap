//
// https://docs.coincap.io
//
// https://api.coincap.io/v2
//

package coincap

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"unsafe"
)

// API
var (
	Assets    assets
	Rates     rates
	Exchanges exchanges
	Markets   markets
	Candles   candles
)

// Currency contains currency id and symbol.
type Currency struct {
	ID     string
	Symbol string
}

// Some currencies
var (
	USD = Currency{"united-states-dollar", "USD"}
	BTC = Currency{"bitcoin", "BTC"}
	ETH = Currency{"ethereum", "ETH"}
)

// Client is a default client that is used to execute requests.
var Client http.Client

var defaultHeader = http.Header{"Accept-Encoding": {"gzip"}}

func request(dataValue interface{}, endPoint string, query url.Values) Timestamp {
	resp, err := Client.Do(&http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.coincap.io",
			Path:     "/v2/" + endPoint,
			RawQuery: query.Encode(),
		},
		Header: defaultHeader,
	})
	if err != nil {
		panic("unexcepted client do error: " + err.Error())
	}

	var body = resp.Body
	defer body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			panic("coincap invalid gzip: " + err.Error())
		}
	}

	// response is coincap normal response.
	var response struct {
		Data      json.RawMessage `json:"data"`
		Timestamp Timestamp       `json:"timestamp"`
	}
	if resp.StatusCode != 200 ||
		resp.Header.Get("Content-Type") != "application/json" ||
		json.NewDecoder(body).Decode(&response) != nil ||
		json.Unmarshal(response.Data, dataValue) != nil {
		bodyBytes, _ := ioutil.ReadAll(body)
		panic("unexcepted coincap response(code " + strconv.Itoa(resp.StatusCode) + "): " + "\n" + string(bodyBytes))
	}

	return response.Timestamp
}

// Int is an int64 with unmarshal.
//
// CoinCap returns numbers as strings.
type Int struct {
	Val int64
}

// UnmarshalJSON parses the JSON-encoded data and stores the result.
func (i *Int) UnmarshalJSON(b []byte) error {
	s := b2s(b)
	var err error
	if s != "null" {
		i.Val, err = strconv.ParseInt(s[1:len(s)-1], 10, 64)
	}
	return err
}

// Float is a float64 with unmarshal.
//
// CoinCap returns numbers as strings.
type Float struct {
	Val float64
}

// UnmarshalJSON parses the JSON-encoded data and stores the result.
func (f *Float) UnmarshalJSON(b []byte) error {
	s := b2s(b)
	var err error
	if s != "null" {
		f.Val, err = strconv.ParseFloat(s[1:len(s)-1], 64)
	}
	return err
}

// TrimParams contains limit and offset parameters.
type TrimParams struct {
	Limit  uint // maximum number of results.
	Offset uint // skip the first N entries of the result set.
}

func (p *TrimParams) set(v *url.Values) {
	if p == nil {
		return
	}
	if p.Limit != 0 {
		if p.Limit > 2000 {
			p.Limit = 2000
		}
		v.Set("limit", utoa(p.Limit))
	}
	v.Set("offset", utoa(p.Offset))
}

// b2s converts bytes slice to string without allocation.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// utoa returns the string representation of num.
func utoa(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}
