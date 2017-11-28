package pingdom

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// test client
	client = NewClient("fake_email@example.com", "12345", "my_api_key")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("user", "password", "key")
	if c.client != http.DefaultClient {
		t.Errorf("NewClient client = %v, want http.DefaultClient", c.client)
	}

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}

	if c.Checks == nil {
		t.Errorf("NewClient Checks service is nil")
	}
}

func TestNewMultiUserClient(t *testing.T) {
	c := NewMultiUserClient("user", "password", "key", "account_email")
	if c.AccountEmail == "" {
		t.Errorf("NewMultiUserClient failed to set AccountEmail property")
	}
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "/checks", nil)
	if err != nil {
		t.Errorf("NewRequest returned error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("NewRequest Method returned %+v, want %+v", req.Method, "GET")
	}

	if req.URL.String() != client.BaseURL.String()+"/checks" {
		t.Errorf("NewRequest URL returned %+v, want %+v", req.URL.String(), client.BaseURL.String()+"/checks")
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestValidateResponse(t *testing.T) {
	valid := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}

	if err := validateResponse(valid); err != nil {
		t.Errorf("ValidateResponse with valid response returned error %+v", err)
	}

	invalid := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{
			"error" : {
				"statuscode": 400,
				"statusdesc": "Bad Request",
				"errormessage": "This is an error"
			}
		}`)),
	}

	want := &PingdomError{400, "Bad Request", "This is an error"}
	if err := validateResponse(invalid); !reflect.DeepEqual(err, want) {
		t.Errorf("ValidateResponse with invalid response returned %+v, want %+v", err, want)
	}

}
