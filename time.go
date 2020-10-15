package coincap

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// IntervalParams contains interval and time span parameters.
type IntervalParams struct {
	Interval Interval  // point-in-time interval.
	Start    time.Time // start time.
	End      time.Time // end time.
}

func (p *IntervalParams) set(v *url.Values, candlesIntervals bool) error {
	if p == nil {
		v.Set("interval", intervals[Hour].str)
		return nil
	}

	il := len(intervals) - 1
	if !candlesIntervals {
		il = extraIntervals - 1
	}

	if uint(p.Interval) > uint(il) {
		return errors.New("invalid interval: use Hour, Minute etc")
	}
	v.Set("interval", intervals[p.Interval].str)

	if p.Start.IsZero() && p.End.IsZero() {
		return nil
	}

	if span := p.End.Sub(p.Start); span < 0 || p.Start.IsZero() || p.End.After(time.Now()) {
		return errors.New("invalid time span")
	} else if span < intervals[p.Interval].dur {
		return errors.New("invalid interval: bigger then time span")
	}

	v.Set("start", MakeTimestamp(p.Start))
	v.Set("end", MakeTimestamp(p.End))

	return nil
}

// Timestamp is UNIX time in milliseconds.
type Timestamp int64

// Time converts CoinCap timestamp into local time.
func (t Timestamp) Time() time.Time {
	return time.Unix(0, int64(t)*1e6)
}

// MakeTimestamp converts local time into CoinCap timestamp.
func MakeTimestamp(ltime time.Time) string {
	return strconv.FormatInt(ltime.UnixNano()/1e6, 10)
}

// Interval represents point-in-time intervals for retrieving historical market data
type Interval int

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

const extraIntervals = int(FourHours)

var intervals = [...]struct {
	str string
	dur time.Duration
}{
	Hour:           {"h1", time.Hour},
	Minute:         {"m1", time.Minute},
	FiveMinutes:    {"m5", 5 * time.Minute},
	FifteenMinutes: {"m15", 15 * time.Minute},
	ThirtyMinutes:  {"m30", 30 * time.Minute},
	TwoHours:       {"h2", 2 * time.Hour},
	SixHours:       {"h6", 6 * time.Hour},
	TwelveHours:    {"h12", 12 * time.Hour},
	Day:            {"d1", 24 * time.Hour},
	FourHours:      {"h4", 4 * time.Hour},
	EightHours:     {"h8", 8 * time.Hour},
	Week:           {"w1", 7 * 24 * time.Hour},
}
