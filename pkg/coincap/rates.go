package coincap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Rate contains the exchange rate of a given asset in terms of USD as well as
// common identifiers for the asset in question and whether or not it is a fiat currency
type Rate struct {
	ID             string `json:"id"`             // unique identifier for asset or fiat
	Symbol         string `json:"symbol"`         // most common symbol used to identify asset or fiat
	CurrencySymbol string `json:"currencySymbol"` // currency symbol if available
	RateUSD        string `json:"rateUsd"`        // rate conversion to USD
	Type           string `json:"type"`           // type of currency - fiat or crypto
}

// Rates returns currency rates standardized in USD.
// Fiat rates are sourced from OpenExchangeRates.org
func (c *Client) Rates() ([]*Rate, *Timestamp, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/rates", nil)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var rates []*Rate
	if err := json.Unmarshal(*ccResp.Data, &rates); err != nil {
		return nil, nil, err
	}

	return rates, ccResp.Timestamp, nil
}

// RateByID returns the USD rate for the given asset identifier
func (c *Client) RateByID(id string) (*Rate, *Timestamp, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/rates/"+id, nil)
	fmt.Println(req)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var rate Rate
	json.Unmarshal(*ccResp.Data, &rate)

	return &rate, ccResp.Timestamp, nil
}
