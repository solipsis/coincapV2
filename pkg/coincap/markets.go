package coincap

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// MarketsRequest contains the paramters you can use to tailor a request for market data from the /markets endpoint
type MarketsRequest struct {
	ExchangeID  string `json:"exchangeId,omitempty"`  // search by unique exchange ID
	BaseSymbol  string `json:"baseSymbol,omitempty"`  // return all results with this base symbol
	BaseID      string `json:"baseId,omitempty"`      // return all results with this base id
	QuoteSymbol string `json:"quoteSymbol,omitempty"` // return all results with this quote symbol
	QuoteID     string `json:"quoteId,omitempty"`     // return all results with this quote ID
	AssetSymbol string `json:"assetSymbol,omitempty"` // return all results with this asset symbol
	AssetID     string `json:"assetID,omitempty"`     // return all results with this asset ID
	Limit       int    `json:"limit,omitempty"`       // limit number of returned results (Max: 2000)
	Offset      int    `json:"offset,omitempty"`      // skip the first N entries of the result set
}

// Market contains the market data response from the api
type Market struct {
	ExchangeID            string    `json:"exchangeId"`            // unique identifier for exchange
	Rank                  string    `json:"rank"`                  // rank in terms of volume transacted in this market
	BaseSymbol            string    `json:"baseSymbol"`            // most common symbol used to identify this asset
	BaseID                string    `json:"baseId"`                // unique identifier for this asset. base is the asset purchased
	QuoteSymbol           string    `json:"quoteSymbol"`           // most common symbol used to identify this asset
	QuoteID               string    `json:"quoteId"`               // unique identifier for thisasset. quote is the asset used to purchase base
	PriceQuote            string    `json:"priceQuote"`            // amount of quote asset traded for 1 unit of base asset
	PriceUsd              string    `json:"priceUsd"`              // quote price translated to USD
	VolumeUsd24Hr         string    `json:"volumeUsd24Hr"`         // volume transacted in this market in the last 24 hours
	PercentExchangeVolume string    `json:"percentExchangeVolume"` // amount of daily volume this market transacts compared to others on this exchange
	TradesCount24Hr       string    `json:"tradesCount24Hr"`       // number of trades on this market in the last 24 hours
	Updated               Timestamp `json:"updated"`               // last time information was received from this market
}

// Markets requests market data for all markets matching the criteria set in the MarketRequest params.
// For historical data on markets use the Candles() endpoint.
// GET /markets
func (c *Client) Markets(reqParams *MarketsRequest) ([]*Market, *Timestamp, error) {

	// Prepare the query
	req, err := http.NewRequest("GET", baseURL+"markets", nil)
	if err != nil {
		return nil, nil, err
	}

	// encode optional parameters
	params := req.URL.Query()
	if reqParams.ExchangeID != "" {
		params.Add("exchange", reqParams.ExchangeID)
	}
	if reqParams.BaseSymbol != "" {
		params.Add("baseSymbol", reqParams.BaseSymbol)
	}
	if reqParams.BaseID != "" {
		params.Add("baseId", reqParams.BaseID)
	}
	if reqParams.QuoteSymbol != "" {
		params.Add("quoteSymbol", reqParams.QuoteSymbol)
	}
	if reqParams.QuoteID != "" {
		params.Add("quoteId", reqParams.QuoteID)
	}
	if reqParams.AssetSymbol != "" {
		params.Add("AssetSymbol", reqParams.AssetSymbol)
	}
	if reqParams.AssetID != "" {
		params.Add("assetId", reqParams.AssetID)
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
	var markets []*Market
	json.Unmarshal(*ccResp.Data, &markets)

	return markets, &ccResp.Timestamp, nil
}
