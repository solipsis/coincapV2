package coincap

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestExchanges(t *testing.T) {

	teardown := setup()
	defer teardown()

	r.HandleFunc("/exchanges", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("exchange.json"))
	})

	exchanges, _, err := client.Exchanges()
	if err != nil {
		t.Fatal(err)
	}
	if len(exchanges) < 0 {
		t.Fatalf("No Exchanges returned")
	}
	got := exchanges[0]
	ts := Timestamp{
		Time: time.Unix(0, 1536336916333*1e6),
	}

	expected := Exchange{
		ID:                 "binance",
		Name:               "Binance",
		Rank:               "1",
		PercentTotalVolume: "16.903027981466749702000000000000000000",
		VolumeUSD:          "1034850514.6425770861221546",
		TradingPairs:       "375",
		Socket:             true,
		Updated:            ts,
	}
	if *got != expected {
		t.Errorf("Expected %v, Got %v", expected, *got)
	}
}

func TestExchangeBadURL(t *testing.T) {
	teardown := setup()
	defer teardown()

	client.SetBaseURL("bad.url")
	_, _, err := client.Exchanges()
	if err == nil {
		t.Errorf("Expected call to fail because of invalid URL")
	}

	_, _, err = client.ExchangeByID("binance")
	if err == nil {
		t.Errorf("Expected call to fail because of invalid URL")
	}

}

func TestExchangeByID(t *testing.T) {

	teardown := setup()
	defer teardown()

	r.HandleFunc("/exchanges/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("exchangeByID.json"))
	})

	exchange, _, err := client.ExchangeByID("gdax")
	if err != nil {
		t.Fatal(err)
	}

	if exchange.ID != "gdax" {
		t.Errorf("Expected exchange ID to be %s, but was %s", "gdax", exchange.ID)
	}

}
