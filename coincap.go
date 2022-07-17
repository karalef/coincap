package coincap

import (
	"encoding/json"
	"errors"
	"io"
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

var DefaultClient = NewClient(nil)

// NewClient creates new CoinCap client.
func NewClient(httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
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

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		b, _ := io.ReadAll(io.MultiReader(dec.Buffered(), resp.Body))
		return 0, errors.New("unexpected CoinCap response:\n" + string(b))
	}

	return response.Timestamp, json.Unmarshal(response.Data, dst)
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

// Interval type.
type Interval uint8

// Valid Intervals for historical market data
// Used when requesting Asset History and Candles
const (
	Hour Interval = iota
	Minute
	FiveMinutes
	FifteenMinutes
	ThirtyMinutes
	TwoHours
	SixHours
	TwelveHours
	Day
	FourHours
	EightHours
	Week
)

func (i Interval) data() (string, time.Duration, bool) {
	switch i {
	case Hour:
		return "h1", time.Hour, false
	case Minute:
		return "m1", time.Minute, false
	case FiveMinutes:
		return "m5", 5 * time.Minute, false
	case FifteenMinutes:
		return "m15", 15 * time.Minute, false
	case ThirtyMinutes:
		return "m30", 30 * time.Minute, false
	case TwoHours:
		return "h2", 2 * time.Hour, false
	case SixHours:
		return "h6", 6 * time.Hour, false
	case TwelveHours:
		return "h12", 12 * time.Hour, false
	case Day:
		return "d1", 24 * time.Hour, false
	case FourHours:
		return "h4", 4 * time.Hour, true
	case EightHours:
		return "h8", 8 * time.Hour, true
	case Week:
		return "w1", 7 * 24 * time.Hour, true
	}
	return "", 0, false
}

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

func (p *TrimParams) setTo(v *url.Values) {
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

// IntervalParams contains interval and time span parameters.
type IntervalParams struct {
	Interval Interval  // point-in-time interval.
	Start    time.Time // start time.
	End      time.Time // end time.
}

func (p *IntervalParams) setTo(v *url.Values, candles bool) error {
	if p == nil {
		p = &IntervalParams{}
	}

	str, dur, ext := p.Interval.data()
	if ext && !candles || str == "" {
		return ErrInvalidInterval
	}

	v.Set("interval", str)

	if p.Start.IsZero() && p.End.IsZero() {
		return nil
	}

	if span := p.End.Sub(p.Start); span < 0 || p.Start.IsZero() || p.End.After(time.Now()) {
		return ErrInvalidTimeSpan
	} else if span < dur {
		return ErrIntervalBigger
	}

	v.Set("start", MakeTimestamp(p.Start).String())
	v.Set("end", MakeTimestamp(p.End).String())

	return nil
}

// MakeTimestamp converts human-readable time into CoinCap timestamp.
func MakeTimestamp(humanTime time.Time) Timestamp {
	return Timestamp(humanTime.UnixNano() / 1e6)
}
