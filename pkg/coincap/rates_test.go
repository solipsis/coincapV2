package coincap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var (
	r      *mux.Router
	server *httptest.Server
	client *Client
)

func setup() func() {
	r = mux.NewRouter()
	server = httptest.NewServer(r)

	client = NewClient(nil)
	client.SetBaseURL(server.URL)

	return func() {
		server.Close()
	}
}

// load test data
func fixture(path string) string {
	f, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic(err)
	}
	return string(f)
}

func TestRates(t *testing.T) {
	teardown := setup()
	defer teardown()

	r.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("rates.json"))
	})

	rates, _, err := client.Rates()
	if err != nil {
		t.Fatal(err)
	}
	if len(rates) < 0 {
		t.Fatalf("No rates returned")
	}
	got := rates[0]
	expected := Rate{
		ID:             "romanian-leu",
		Symbol:         "RON",
		CurrencySymbol: "lei",
		Type:           "fiat",
		RateUSD:        "0.2505529076289101",
	}
	if *got != expected {
		t.Errorf("Expected %s, Got %s", expected, got)
	}
}

func TestRatesBadURL(t *testing.T) {
	teardown := setup()
	defer teardown()

	client.SetBaseURL("bad.fail:")
	_, _, err := client.Rates()
	if err == nil {
		t.Errorf("Expected client to fail on Rates() because of bad base path")
	}
	_, _, err = client.RateByID("bitcoin")
	if err == nil {
		t.Errorf("Expected client to fail on RateByID() because of bad base path")
	}

}

func TestRateByID(t *testing.T) {
	teardown := setup()
	defer teardown()

	r.HandleFunc("/rates/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("ratesByID.json"))
	})

	rate, _, err := client.RateByID("bitcoin")
	if err != nil {
		t.Fatal(err)
	}
	expected := Rate{
		ID:             "bitcoin",
		Symbol:         "BTC",
		CurrencySymbol: "₿",
		Type:           "crypto",
		RateUSD:        "6460.9771089680171173",
	}
	if *rate != expected {
		t.Errorf("Expected %s, Got %s", expected, rate)
	}
}

func TestRatesMalformed(t *testing.T) {
	teardown := setup()
	defer teardown()

	r.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("rates_malformed.json"))
	})

	rates, _, err := client.Rates()
	if err == nil {
		t.Errorf("Expected malformed json %s", rates)
	}
}

func TestRatesNoTimestamp(t *testing.T) {
	teardown := setup()
	defer teardown()

	r.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("rates_no_timestamp.json"))
	})

	_, timestamp, err := client.Rates()
	if err == nil {
		t.Errorf("Expected error due to missing timestamp json %s", timestamp)
	}
}
