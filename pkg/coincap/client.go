package coincap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient returns a new client for interacting with the CoinCap API
// If no httpClient is passed it will use http.DefaultClient
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
}

// SetBaseURL allows the setting of custom api base paths
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// Every coincap response has a top level entry called data
// and a unix timestamp in milliseconds
type coincapResp struct {
	Data      *json.RawMessage `json:"data"`
	Timestamp Timestamp        `json:"timestamp"`
}

// fetchAndParse returns the json below the top level "data" key
// returned by the coincap api
func (c *Client) fetchAndParse(req *http.Request) (*coincapResp, error) {

	// make request to the api and read the response
	resp, err := c.httpClient.Do(req)
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
	fmt.Println("Body:", string(body))

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
