# GoLang CoinCap API

Library for interaction with [CoinCap API 2.0](https://docs.coincap.io/). Supports all methods and types of responses including WebSockets.

## How to install

```bash
go get -u github.com/karalef/coincap
```

## Testing

```bash
go test .
```

## Usage

```go
import "github.com/karalef/coincap"

...

client := coincap.DefaultClient
```

### Custom client

```go
import "github.com/gorilla/websocket"

wsDialer := websocket.Dialer{
    ...
}

client := coincap.New(&http.Client{}, &wsDialer)
```

## Examples

### Get Asset data

```go
asset, timestamp, err := client.AssetById("bitcoin")
```

### Get historical data

```go
params := coincap.CandlesRequest{
    ExchangeID: "binance",
    BaseID:     "ethereum",
    QuoteID:    "bitcoin",
}

interval := coincap.IntervalParams{
    Interval: coincap.Hour,
    Start:    time.Now().Add(-time.Day*5),
    End:      time.Now(),
}

candles, timestamp, err := client.Candles(params, &interval, nil)
```

### WebSockets

```go
stream, err := client.Trades("binance")
if err != nil {
    ...
}

ch := stream.DataChannel()
for {
    trade, ok := <-ch
    if !ok {
        err := stream.Err()
        ...
    }
    ...
}
stream.Close()
```


