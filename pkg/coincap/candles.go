package coincap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// CandlesRequest contains the parameters you can use to customize a request for candle data from the /candles endpoint
type CandlesRequest struct {
	ExchangeID string   `json:"exchangeId"`       // search by unique exchange ID
	BaseID     string   `json:"baseId"`           // return all results with this base id
	QuoteID    string   `json:"quoteId"`          // return all results with this quote ID
	Interval   Interval `json:"interval"`         // candle interval
	Start      int      `json:"start,omitempty"`  // start time in unix milliseconds TODO: I should probably use time.Time or Timestamp here
	End        int      `json:"end,omitempty"`    // end time in unix milliseconds TODO: same as above
	Limit      int      `json:"limit,omitempty"`  // limit number of returned results (Max: 2000)
	Offset     int      `json:"offset,omitempty"` // skip the first N entries of the result set
}

// Candle represets historic market performance for an asset over a given timeframe
type Candle struct {
	Open   string    `json:"open"`   // the price (quote) at which the first transaction was completed in a given time period
	High   string    `json:"high"`   // the top price (quote) at which the base was traded during the time period
	Low    string    `json:"low"`    // the bottom price (quote) at which the base was traded during the time period
	Close  string    `json:"close"`  // the price (quote) at which the last transaction was completed in a given time period
	Volume string    `json:"volume"` // the amount of base asset traded in the given time period
	Period Timestamp `json:"period"` // timestamp for starting of that time period
}

// Candles returns all the market candle data for the provided exchange and parameters
// The fields ExchangeID, BaseID, QuoteID, and Interval are required by the API
func (c *Client) Candles(reqParams *CandlesRequest) ([]*Candle, *Timestamp, error) {

	// check required params
	var err error
	if reqParams.ExchangeID == "" {
		err = fmt.Errorf("ExchangeID is required")
	} else if reqParams.BaseID == "" {
		err = fmt.Errorf("BaseID is required")
	} else if reqParams.QuoteID == "" {
		err = fmt.Errorf("QuoteID is required")
	} else if string(reqParams.Interval) == "" {
		err = fmt.Errorf("Interval is required")
	}
	if err != nil {
		return nil, nil, err
	}

	// Prepare the query
	req, err := http.NewRequest("GET", baseURL+"candles", nil)
	if err != nil {
		return nil, nil, err
	}

	// encode parameters
	params := req.URL.Query()
	params.Add("exchange", reqParams.ExchangeID)
	params.Add("baseId", reqParams.BaseID)
	params.Add("quoteId", reqParams.QuoteID)
	params.Add("interval", string(reqParams.Interval))
	if reqParams.Start > 0 {
		params.Add("start", strconv.Itoa(reqParams.Start))
	}
	if reqParams.End > 0 {
		params.Add("end", strconv.Itoa(reqParams.End))
	}
	if reqParams.Limit > 0 {
		params.Add("limit", strconv.Itoa(reqParams.Limit))
	}
	if reqParams.Offset > 0 {
		params.Add("offset", strconv.Itoa(reqParams.Offset))
	}
	req.URL.RawQuery = params.Encode()
	fmt.Println(req.URL.RawQuery)

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var candles []*Candle
	json.Unmarshal(*ccResp.Data, &candles)

	return candles, &ccResp.Timestamp, nil
}
