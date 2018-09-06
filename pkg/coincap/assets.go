package coincap

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// AssetsRequest contains the paramaters for modifying a query to
// the "/assets" endpoint. Search can be a symbol (BTC) or an asset id (bitcoin)
type AssetsRequest struct {
	Search string `json:"search,omitempty"` // search by asset id (bitcoin) or symbol (BTC)
	Limit  int    `json:"limit,omitempty"`  // limit number of returned results (Max: 2000)
	Offset int    `json:"offset,omitempty"` // skip the first N entries of the result set
}

// Asset contains various information about a given CoinCap asset such as Bitcoin
type Asset struct {
	ID                string `json:"id"`                // unique identifier for asset
	Rank              string `json:"rank"`              // rank in terms of the asset's market cap
	Symbol            string `json:"symbol"`            // common symbol to identify the asset
	Name              string `json:"name"`              // proper name for asset
	Supply            string `json:"supply"`            // available supply for trading
	MaxSupply         string `json:"maxSupply"`         // total quantity of asset issued
	MarketCapUsd      string `json:"marketCapUsd"`      // supply x price
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`     // quantity of trading volume in USD over last 24 hours
	PriceUsd          string `json:"priceUsd"`          // volume weighted price of the asset in USD
	ChangePercent24Hr string `json:"changePercent24Hr"` // percent change in value in the last 24 hours
	Vwap24Hr          string `json:"vwap24Hr"`          // Volume Weighted Average Price in the last 24 hours
}

// Assets returns a list of CoinCap Asset entries filtered by the request's
// search criteria and a timestamp
func (c *Client) Assets(reqParams *AssetsRequest) ([]Asset, *Timestamp, error) {

	// Prepare the query and encode optional parameters
	req, err := http.NewRequest("GET", baseURL+"assets", nil)
	if err != nil {
		return nil, nil, err
	}
	params := req.URL.Query()
	params.Add("search", reqParams.Search)
	if reqParams.Limit > 0 {
		params.Add("limit", strconv.Itoa(reqParams.Limit))
	}
	if reqParams.Offset > 0 {
		params.Add("offset", strconv.Itoa(reqParams.Offset))
	}
	req.URL.RawQuery = params.Encode()

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var assets []Asset
	json.Unmarshal(*ccResp.Data, &assets)

	return assets, ccResp.Timestamp, nil
}

// AssetByID requests an asset by its CoinCap ID
func (c *Client) AssetByID(id string) (Asset, *Timestamp, error) {

	req, err := http.NewRequest("GET", baseURL+"assets/"+id, nil)
	if err != nil {
		return Asset{}, nil, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return Asset{}, nil, err
	}

	// Unmarshal the deferred json from the data field
	var asset Asset
	json.Unmarshal(*ccResp.Data, &asset)

	return asset, ccResp.Timestamp, nil
}

// AssetHistoryRequest contains the paramaters for modifying a query to
// the "/assets/{{id}}/history" endpoint.
type AssetHistoryRequest struct {
	Interval Interval `json:"interval"`         // point-in-time interval. minute, hour, and day. Allowed intervals (m1, m15, h1, d1)
	Start    int      `json:"start,omitempty"`  // start time in unix milliseconds TODO: I should probably use time.Time or Timestamp here
	End      int      `json:"end,omitempty"`    // end time in unix milliseconds TODO: same as above
	Limit    int      `json:"limit,omitempty"`  // maximum number of results to return
	Offset   int      `json:"offset,omitempty"` // skip some of the returned results
}

// AssetHistory contains the USD price of an asset at a given timestamp
type AssetHistory struct {
	PriceUSD string    `json:"priceUsd"` // volume weighted price in USD
	Time     Timestamp `json:"time"`     // Timestamp correlating to the given price
}

// AssetHistoryByID returns USD price history of a given asset.
// If no interval is specified 1 hour (h1) is chosen as the default.
func (c *Client) AssetHistoryByID(id string, reqParams *AssetHistoryRequest) ([]AssetHistory, *Timestamp, error) {

	// Default interval to an hour if none was provided
	if reqParams.Interval == "" {
		reqParams.Interval = Hour
	}

	// Prepare the query
	req, err := http.NewRequest("GET", baseURL+"assets/"+id+"/history", nil)
	if err != nil {
		return nil, nil, err
	}

	// encode optional parameters
	params := req.URL.Query()
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

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var history []AssetHistory
	json.Unmarshal(*ccResp.Data, &history)

	return history, ccResp.Timestamp, nil
}
