package coincap

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//
// https://docs.coincap.io
//
// https://api.coincap.io/v2
//

// NewClient creates new CoinCap client.
func NewClient(httpClient *http.Client) Client {
	return Client{httpClient}
}

// Client gives access to CoinCap API.
type Client struct {
	http *http.Client
}

func (c *Client) request(dst interface{}, endPoint string, query url.Values) (Timestamp, error) {
	resp, err := c.http.Do(&http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.coincap.io",
			Path:     "/v2/" + endPoint,
			RawQuery: query.Encode(),
		},
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("unexpected http error with status: " + resp.Status)
	}

	// response is a CoinCap normal response.
	var response struct {
		Data      json.RawMessage `json:"data"`
		Timestamp Timestamp       `json:"timestamp"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil || json.Unmarshal(response.Data, dst) != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return 0, errors.New("unexpected CoinCap response: \n" + string(bodyBytes))
	}

	return response.Timestamp, nil
}

// Timestamp represents CoinCap timestamp
// (UNIX time in milliseconds).
type Timestamp int64

// Time converts CoinCap timestamp into local time.
func (t Timestamp) Time() time.Time {
	return time.Unix(0, int64(t)*1e6)
}

func (t Timestamp) String() string {
	return strconv.FormatInt(int64(t), 10)
}

// MakeTimestamp converts human-readable time into CoinCap timestamp.
func MakeTimestamp(humanTime time.Time) Timestamp {
	return Timestamp(humanTime.UnixNano() / 1e6)
}

// interval represents point-in-time intervals for retrieving historical market data
type interval struct {
	str string
	dur time.Duration
	ext bool
}

// Valid Intervals for historical market data
// Used when requesting Asset History and Candles
var (
	Hour           = interval{"h1", time.Hour, false}
	Minute         = interval{"m1", time.Minute, false}
	FiveMinutes    = interval{"m5", 5 * time.Minute, false}
	FifteenMinutes = interval{"m15", 15 * time.Minute, false}
	ThirtyMinutes  = interval{"m30", 30 * time.Minute, false}
	TwoHours       = interval{"h2", 2 * time.Hour, false}
	SixHours       = interval{"h6", 6 * time.Hour, false}
	TwelveHours    = interval{"h12", 12 * time.Hour, false}
	Day            = interval{"d1", 24 * time.Hour, false}
	FourHours      = interval{"h4", 4 * time.Hour, true}
	EightHours     = interval{"h8", 8 * time.Hour, true}
	Week           = interval{"w1", 7 * 24 * time.Hour, true}
)

// interval errors.
var (
	ErrInvalidInterval = errors.New("invalid interval")
	ErrInvalidTimeSpan = errors.New("invalid time span")
	ErrIntervalBigger  = errors.New("invalid interval: bigger then time span")
)

func utoa(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}

// TrimParams contains limit and offset parameters.
type TrimParams struct {
	Limit  uint // maximum number of results.
	Offset uint // skip the first N entries of the result set.
}

// IntervalParams contains interval and time span parameters.
type IntervalParams struct {
	Interval interval  // point-in-time interval.
	Start    time.Time // start time.
	End      time.Time // end time.
}

func setTrim(p *TrimParams, v *url.Values) {
	if p == nil {
		return
	}
	if p.Limit != 0 {
		if p.Limit > 2000 {
			p.Limit = 2000
		}
		v.Set("limit", utoa(p.Limit))
	}
	if p.Offset != 0 {
		v.Set("offset", utoa(p.Offset))
	}
}

func setInterval(p *IntervalParams, v *url.Values, candles bool) error {
	if p == nil {
		v.Set("interval", Hour.str)
		return nil
	}

	if !candles && p.Interval.ext {
		return ErrInvalidInterval
	}

	v.Set("interval", p.Interval.str)

	if p.Start.IsZero() && p.End.IsZero() {
		return nil
	}

	if span := p.End.Sub(p.Start); span < 0 || p.Start.IsZero() || p.End.After(time.Now()) {
		return ErrInvalidTimeSpan
	} else if span < p.Interval.dur {
		return ErrIntervalBigger
	}

	v.Set("start", MakeTimestamp(p.Start).String())
	v.Set("end", MakeTimestamp(p.End).String())

	return nil
}
