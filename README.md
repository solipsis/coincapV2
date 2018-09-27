
# coincapV2 #

[![GoDoc](https://godoc.org/github.com/solipsis/coincapV2?status.svg)](https://godoc.org/github.com/solipsis/coincapV2) [![Go Report Card](https://goreportcard.com/badge/github.com/solipsis/coincapV2)](https://goreportcard.com/report/github.com/solipsis/coincapV2) 

client library for interacting with the coincap.io V2 api

## Installation ##

	go get -u github.com/solipsis/coincapV2/...
  

## Usage ##

```go
import "github.com/solipsis/coincapV2/pkg/coincap"
```
GoDoc: https://godoc.org/github.com/solipsis/coincapV2/pkg/coincap

## Official API Docs ##
https://docs.coincap.io/


## Examples ##

### Get Market Data ###
```go
client := coincap.NewClient(nil)

params := &MarketsRequest{
	ExchangeID:  "binance",
	BaseSymbol:  "ETH",
	QuoteSymbol: "BTC",
	Limit:       100,
	Offset:      0,
}
markets, timestamp, err := client.Markets(params)
```

### Get Data for an Asset ###
```go
client := coincap.NewClient(nil)

params := &AssetsRequest{
	Search: "BTC",
	Limit:  4,
	Offset: 1,
}
assets, timestamp, err := client.Assets(params)
```

### Get Historical Data for an Asset ###

```go
client := coincap.NewClient(nil)

// setup the time range
end := time.Now()
start := now.Add(-time.Hour * 2)

params := &coincap.AssetHistoryRequest{
	Interval: coincap.FifteenMinutes,
	Start:    &coincap.Timestamp{Time: start},
	End:      &coincap.Timestamp{Time: end},
}
history, timestamp, err := client.AssetHistoryByID("bitcoin", params)
if err != nil {
	log.Fatal(err)
}
```

### Get Rates of Various Currencies to USD ###

```go
client := coincap.NewClient(nil)

rates, timestamp, err := client.Rates()
if err != nil {
	t.Fatal(err)
}

/*
Example response
{
  "data": [
    {
      "id": "barbadian-dollar",
      "symbol": "BBD",
      "currencySymbol": "$",
      "type": "fiat",
      "rateUsd": "0.5000000000000000"
    },
    {
      "id": "malawian-kwacha",
      "symbol": "MWK",
      "currencySymbol": "MK",
      "type": "fiat",
      "rateUsd": "0.0013750106599420"
    },
    {
      "id": "south-african-rand",
      "symbol": "ZAR",
      "currencySymbol": "R",
      "type": "fiat",
      "rateUsd": "0.0657544075508153"
    }],
}
```

### Get USD Rates by Asset ID ###

```go
client := coincap.NewClient(nil)

rate, timestamp, err := client.RateByID("bitcoin")
```

### Get Information on Exchanges ###

```go
client := coincap.NewClient(nil)

exchanges, timestamp, err := client.Exchanges()
if err != nil {
	t.Fatal(err)
}
```

### Get Exchange Information by ID ###

```go
client := coincap.NewClient(nil)

exchange, timestamp, err := client.ExchangeByID("gdax")
if err != nil {
	t.Fatal(err)
}
```

### Get Candle Data ###
```go
client := coincap.NewClient(nil)

params := &CandlesRequest{
	ExchangeID: "poloniex",
	BaseID:     "ethereum",
	QuoteID:    "bitcoin",
	Limit:      100,
	Offset:     1,
	Interval:   coincap.FiveMinutes,
}
candles, timestamp, err := client.Candles(params)
if err != nil {
	t.Fatal(err)
}
```

## Contributing ##
Contributions and pull requests welcome
