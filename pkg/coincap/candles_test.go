package coincap

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCandles(t *testing.T) {
	teardown := setup()
	defer teardown()
	r.HandleFunc("/candles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("candles.json"))
	})

	req := CandlesRequest{
		ExchangeID: "poloniex",
		BaseID:     "ethereum",
		QuoteID:    "bitcoin",
		Limit:      100,
		Offset:     1,
		Interval:   FiveMinutes,
	}
	_, _, err := client.Candles(&req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFailRequiredParams(t *testing.T) {
	teardown := setup()
	defer teardown()

	// exchangeID
	req := CandlesRequest{
		QuoteID:  "bitcoin",
		Limit:    100,
		Offset:   1,
		Interval: FiveMinutes,
	}
	_, _, err := client.Candles(&req)
	if err == nil {
		t.Errorf("Expected client to fail because all required paramters were not provided")
	}

	// quoteID
	req = CandlesRequest{
		ExchangeID: "poloniex",
		BaseID:     "ethereum",
		Limit:      100,
		Offset:     1,
	}
	_, _, err = client.Candles(&req)
	if err == nil {
		t.Errorf("Expected client to fail because all required paramters were not provided")
	}

	// BaseID
	req = CandlesRequest{
		ExchangeID: "poloniex",
		QuoteID:    "bitcoin",
		Interval:   FiveMinutes,
	}
	_, _, err = client.Candles(&req)
	if err == nil {
		t.Errorf("Expected client to fail because all required paramters were not provided")
	}

	// interval
	req = CandlesRequest{
		ExchangeID: "poloniex",
		QuoteID:    "bitcoin",
		BaseID:     "ethereum",
	}
	_, _, err = client.Candles(&req)
	if err == nil {
		t.Errorf("Expected client to fail because all required paramters were not provided")
	}

}
