package coincap

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMarkets(t *testing.T) {

	teardown := setup()
	defer teardown()

	r.HandleFunc("/markets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("markets.json"))
	})

	params := &MarketsRequest{
		ExchangeID:  "binance",
		BaseSymbol:  "ETH",
		QuoteSymbol: "BTC",
		Limit:       100,
		Offset:      0,
	}

	markets, _, err := client.Markets(params)
	if err != nil {
		t.Fatal(err)
	}
	if len(markets) < 0 {
		t.Fatalf("No Markets returned")
	}
	got := markets[0]
	if got.ExchangeID != "binance" {
		t.Errorf("Expected exchangeID to be binance but was %s", got.ExchangeID)
	}
}

func TestMarketsLive(t *testing.T) {
	// Hit the actual API
	client := NewClient(nil)

	params := &MarketsRequest{
		BaseID:      "ethereum",
		QuoteID:     "dogecoin",
		AssetSymbol: "dogecoin",
	}
	markets, _, err := client.Markets(params)
	if err != nil {
		t.Fatal(err)
	}

	if len(markets) == 0 {
		t.Errorf("No markets were returned")
	}
}
