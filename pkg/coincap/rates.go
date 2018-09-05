package coincap

import (
	"encoding/json"
	"net/http"
)

// Rate contains the exchange rate of a given asset in terms of USD as well as
// common identifiers for the asset in question and whether or not it is a fiat currency
type Rate struct {
	ID      string `json:"id"`      // unique identifier for asset or fiat
	Symbol  string `json:"symbol"`  // most common symbol used to identify asset or fiat
	RateUSD string `json:"rateUsd"` // rate conversion to USD
	Type    string `json:"type"`    // type of currency - fiat or crypto
}

// Rates returns currency rates standardized in USD.
// Fiat rates are sourced from OpenExchangeRates.org
func (c *Client) Rates() ([]Rate, Timestamp, error) {
	req, err := http.NewRequest("GET", baseURL+"rates", nil)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var rates []Rate
	json.Unmarshal(*ccResp.Data, &rates)

	return rates, ccResp.Timestamp, nil
}

// RateByID returns the USD rate for the given asset identifier
func (c *Client) RateByID(id string) (Rate, Timestamp, error) {
	req, err := http.NewRequest("GET", baseURL+"rates/"+id, nil)
	if err != nil {
		return Rate{}, Timestamp{}, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return Rate{}, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var rate Rate
	json.Unmarshal(*ccResp.Data, &rate)

	return rate, ccResp.Timestamp, nil
}
