package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	ass, timestamp, err := Assets(AssetsRequest{})
	if false {
		fmt.Println(ass)
	}
	fmt.Println(pretty(ass))
	fmt.Println("timestamp:", timestamp)
	fmt.Println("err:", err)

}

type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	fmt.Println("In Unmarshal:", string(b))

	// CoinCap timestamp is unix milliseconds
	m, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	// Convert from milleseconds to nanoseconds
	t.Time = time.Unix(0, m*1e6)
	fmt.Println("Tiiiimme", t.Time)
	return nil
}

type cqResponse struct {
	Data      *json.RawMessage `json:"data"`
	Timestamp Timestamp        `json:"timestamp"`
}

type AssetsRequest struct {
	Search string `json:"search,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

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

// parseData returns the json below the top level "data" key
// returned by the coincap api
func parseData(resp *http.Response) (*json.RawMessage, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]*json.RawMessage)

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	return data["data"], nil
}

// JSON pretty print
func pretty(v interface{}) string {
	str, _ := json.MarshalIndent(v, "", "    ")
	return string(str)
}

func Assets(req AssetsRequest) ([]Asset, Timestamp, error) {
	resp, err := http.Get(baseURL + "assets")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Timestamp{}, nil
	}
	fmt.Println(string(body))

	cqResp := cqResponse{}
	json.Unmarshal(body, &cqResp)

	var assets []Asset

	fmt.Println(cqResp.Timestamp)
	json.Unmarshal(*cqResp.Data, &assets)

	return assets, cqResp.Timestamp, nil
}
