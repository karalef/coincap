package coincap

import (
	"net/url"
	"testing"
	"time"
)

var now = time.Now()

func TestIntervalParams(t *testing.T) {
	val := url.Values{}

	// normal params
	var p = IntervalParams{
		Interval: Minute,
		Start:    now.Add(-1 * time.Hour),
		End:      now,
	}
	err := p.set(&val, false)
	if err != nil {
		t.Fatal(err)
	}

	// nil params
	_ = (*IntervalParams)(nil).set(&val, false) // SIGSEGV with no error
}

func TestInvalidInterval(t *testing.T) {
	val := url.Values{}

	p := IntervalParams{
		Interval: 16, // unexisting interval
	}
	err := p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with unexisting interval")
	}

	p = IntervalParams{
		Interval: Week, // Week is not avaible for assets history
	}
	err = p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with unavaible interval")
	}

	p = IntervalParams{
		Interval: Hour,
		Start:    now.Add(-time.Minute), // interval > time span
		End:      now,
	}
	err = p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with interval > time span")
	}
}

func TestInvalidSpan(t *testing.T) {
	val := url.Values{}

	// no end
	p := IntervalParams{
		Start: now,
	}
	err := p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with no end")
	}

	// no start
	p = IntervalParams{
		End: now,
	}
	err = p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with no start")
	}

	p = IntervalParams{
		Start: now.Add(time.Hour), // start > end
		End:   now,
	}
	p.Start = now.Add(time.Hour)
	err = p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with invalid time span(start > end)")
	}

	p = IntervalParams{
		Start: now.Add(-1 * time.Hour),
		End:   now.Add(time.Hour), // start > end
	}
	err = p.set(&val, false)
	if err == nil {
		t.Fatal("NO error with invalid time span(end > now)")
	}
}
