package pingdom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.pingdom.com/api/3.1"
)

// Client represents a client to the Pingdom API.
type Client struct {
	APIToken     string
	BaseURL      *url.URL
	client       *http.Client
	Checks       *CheckService
	Maintenances *MaintenanceService
	Probes       *ProbeService
}

// ClientConfig represents a configuration for a pingdom client.
type ClientConfig struct {
	APIToken   string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClientWithConfig returns a Pingdom client.
func NewClientWithConfig(config ClientConfig) (*Client, error) {
	var baseURL *url.URL
	var err error
	if config.BaseURL != "" {
		baseURL, err = url.Parse(config.BaseURL)
	} else {
		baseURL, err = url.Parse(defaultBaseURL)
	}
	if err != nil {
		return nil, err
	}

	c := &Client{
		APIToken: config.APIToken,
		BaseURL:  baseURL,
	}

	if config.HTTPClient != nil {
		c.client = config.HTTPClient
	} else {
		c.client = http.DefaultClient
	}

	c.Checks = &CheckService{client: c}
	c.Maintenances = &MaintenanceService{client: c}
	c.Probes = &ProbeService{client: c}
	return c, nil
}

// NewRequest makes a new HTTP Request.  The method param should be an HTTP method in
// all caps such as GET, POST, PUT, DELETE.  The rsc param should correspond with
// a restful resource.  Params can be passed in as a map of strings
// Usually users of the client can use one of the convenience methods such as
// ListChecks, etc but this method is provided to allow for making other
// API calls that might not be built in.
func (pc *Client) NewRequest(method string, rsc string, params map[string]string) (*http.Request, error) {
	baseURL, err := url.Parse(pc.BaseURL.String() + rsc)
	if err != nil {
		return nil, err
	}

	if params != nil {
		ps := url.Values{}
		for k, v := range params {
			ps.Set(k, v)
		}
		baseURL.RawQuery = ps.Encode()
	}

	req, err := http.NewRequest(method, baseURL.String(), nil)
	req.Header.Add("Authorization", "Bearer "+pc.APIToken)
	return req, err
}

// Do makes an HTTP request and will unmarshal the JSON response in to the
// passed in interface.  If the HTTP response is outside of the 2xx range the
// response will be returned along with the error.
func (pc *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return resp, err
	}

	err = decodeResponse(resp, v)
	return resp, err

}

func decodeResponse(r *http.Response, v interface{}) error {
	if v == nil {
		return fmt.Errorf("nil interface provided to decodeResponse")
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	err := json.Unmarshal([]byte(bodyString), &v)
	return err
}

// Takes an HTTP response and determines whether it was successful.
// Returns nil if the HTTP status code is within the 2xx range.  Returns
// an error otherwise.
func validateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	m := &errorJSONResponse{}
	err := json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return err
	}

	return m.Error
}
