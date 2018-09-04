package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	/*
		ass, timestamp, err := Assets(AssetsRequest{Search: "bitcoin", Limit: 3})
		if false {
			fmt.Println(ass)
		}
		fmt.Println(pretty(ass))
		fmt.Println("timestamp:", timestamp)
		fmt.Println("err:", err)
	*/
	//ass, timestamp, err := AssetByID("dash")
	//hist, timestamp, err := AssetHistoryByID("dash", &AssetHistoryRequest{})
	//rate, timestamp, err := RateByID("dogecoin")
	//fmt.Println(pretty(ass))
	//exchanges, timestamp, err := Exchanges()
	_, _, err := ExchangeByID("gdax2")
	//fmt.Println(pretty(exchange))
	//fmt.Println(timestamp)
	fmt.Println(err)

}

// Timestamp is wrapper around time.Time with custom unmarshaling behaviour
// specific to the format returned by the CoinCap API
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
// Custom unmarshaller to handle that the timestamp is not in a standard format
func (t *Timestamp) UnmarshalJSON(b []byte) error {

	// CoinCap timestamp is unix milliseconds
	m, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	// Convert from milleseconds to nanoseconds
	t.Time = time.Unix(0, m*1e6)
	return nil
}

// Every coincap response has a top level entry called data
// and a unix timestamp in milliseconds
type coincapResp struct {
	Data      *json.RawMessage `json:"data"`
	Timestamp Timestamp        `json:"timestamp"`
}

// AssetsRequest contains the paramaters for modifying a query to
// the "/assets" endpoint. Search can be a symbol (BTC) or an asset id (bitcoin)
type AssetsRequest struct {
	Search string `json:"search,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// Asset contains various information about a given CoinCap asset such as Bitcoin
type Asset struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
}

// fetchAndParse returns the json below the top level "data" key
// returned by the coincap api
func fetchAndParse(req *http.Request) (*coincapResp, error) {
	client := &http.Client{} // TODO: client should probably live on an API object and these methods on top of that.

	// make request to the api and read the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error received status: %d, with body: %s", resp.StatusCode, string(body))
	}

	// parse the result
	ccResp := new(coincapResp)
	if err := json.Unmarshal(body, ccResp); err != nil {
		return nil, err
	}

	return ccResp, nil
}

// JSON pretty print
func pretty(v interface{}) string {
	str, _ := json.MarshalIndent(v, "", "    ")
	return string(str)
}

// Assets returns a list of CoinCap Asset entries filtered by the request's
// search criteria and a timestamp
func Assets(reqParams *AssetsRequest) ([]Asset, Timestamp, error) {

	// Prepare the query and encode optional parameters
	req, err := http.NewRequest("GET", baseURL+"assets", nil)
	if err != nil {
		return nil, Timestamp{}, err
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
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var assets []Asset
	json.Unmarshal(*ccResp.Data, &assets)

	return assets, ccResp.Timestamp, nil
}

// AssetByID requests an asset by its CoinCap ID
func AssetByID(id string) (Asset, Timestamp, error) {

	req, err := http.NewRequest("GET", baseURL+"assets/"+id, nil)
	if err != nil {
		return Asset{}, Timestamp{}, err
	}

	// make the request
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return Asset{}, Timestamp{}, err
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
	PriceUSD string    `json:"priceUsd"`
	Time     Timestamp `json:"time"`
}

// Interval represents point-in-time intervals for retrieving historical data
type Interval string

const (
	Minute         Interval = "m1"  // 1 minute
	FifteenMinutes Interval = "m15" // 15 minutes
	Hour           Interval = "h1"  // 1 hour
	Day            Interval = "d1"  // 1 day
)

// AssetHistoryByID returns USD price history of a given asset.
// If no interval is specified 1 hour (h1) is chosen as the default.
func AssetHistoryByID(id string, reqParams *AssetHistoryRequest) ([]AssetHistory, Timestamp, error) {

	// Default interval to an hour if none was provided
	if reqParams.Interval == "" {
		reqParams.Interval = Hour
	}

	// Prepare the query
	req, err := http.NewRequest("GET", baseURL+"assets/"+id+"/history", nil)
	if err != nil {
		return nil, Timestamp{}, err
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
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var history []AssetHistory
	json.Unmarshal(*ccResp.Data, &history)

	return history, ccResp.Timestamp, nil
}

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
func Rates() ([]Rate, Timestamp, error) {
	req, err := http.NewRequest("GET", baseURL+"rates", nil)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// make the request
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return nil, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var rates []Rate
	json.Unmarshal(*ccResp.Data, &rates)

	return rates, ccResp.Timestamp, nil
}

// RateByID returns the USD rate for the given asset identifier
func RateByID(id string) (Rate, Timestamp, error) {
	req, err := http.NewRequest("GET", baseURL+"rates/"+id, nil)
	if err != nil {
		return Rate{}, Timestamp{}, err
	}

	// make the request
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return Rate{}, Timestamp{}, err
	}

	// Unmarshal the deferred json from the data field
	var rate Rate
	json.Unmarshal(*ccResp.Data, &rate)

	return rate, ccResp.Timestamp, nil

}

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
func Exchanges() ([]*Exchange, *Timestamp, error) {

	req, err := http.NewRequest("GET", baseURL+"exchanges", nil)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var exchanges []*Exchange
	json.Unmarshal(*ccResp.Data, &exchanges)

	return exchanges, &ccResp.Timestamp, nil
}

// ExchangeByID returns exchange data for an exchange with the given unique ID.
// GET /exchanges/{{id}}
func ExchangeByID(id string) (*Exchange, *Timestamp, error) {

	req, err := http.NewRequest("GET", baseURL+"exchanges/"+id, nil)
	if err != nil {
		return nil, nil, err
	}

	// make the request
	ccResp, err := fetchAndParse(req)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the deferred json from the data field
	var exchange *Exchange
	json.Unmarshal(*ccResp.Data, &exchange)

	return exchange, &ccResp.Timestamp, nil
}
