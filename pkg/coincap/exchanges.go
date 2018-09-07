package coincap

import (
	"encoding/json"
	"net/http"
)

// Exchange contains information about a cryptocurrency exchange. This includes the exchanges
// relative rank, volume, and whether trading sockets are available
type Exchange struct {
	ID                 string    `json:"id"`                 // unique identifier for exchange
	Name               string    `json:"name"`               // proper name of exchange
	Rank               string    `json:"rank"`               // rank in terms of total volume compared to other exchanges
	PercentTotalVolume string    `json:"percentTotalVolume"` // perecent of total daily volume in relation to all exchanges
	VolumeUSD          string    `json:"volumeUSD"`          // daily volume represented in USD
	TradingPairs       string    `json:"tradingPairs"`       // number of trading pairs offered by the exchange
	Socket             bool      `json:"socket"`             // Whether or not a trade socket is available on this exchange
	Updated            Timestamp `json:"updated"`            // Time since information was last updated
}

// Exchanges returns information about all the various exchanges currently tracked by CoinCap.
// GET /exchanges
func (c *Client) Exchanges() ([]*Exchange, *Timestamp, error) {

	req, err := http.NewRequest("GET", c.baseURL+"/exchanges", nil)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var exchanges []*Exchange
	json.Unmarshal(*ccResp.Data, &exchanges)

	return exchanges, ccResp.Timestamp, nil
}

// ExchangeByID returns exchange data for an exchange with the given unique ID.
// GET /exchanges/{{id}}
func (c *Client) ExchangeByID(id string) (*Exchange, *Timestamp, error) {

	req, err := http.NewRequest("GET", c.baseURL+"/exchanges/"+id, nil)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := c.fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var exchange *Exchange
	json.Unmarshal(*ccResp.Data, &exchange)

	return exchange, ccResp.Timestamp, nil
}
