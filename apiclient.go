package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
	"reflect"
)

type API struct {
	baseUrl *url.URL
	Headers map[string]string
	client  *http.Client
}

func New(baseUrl string) *API {
	a := &API{
		client: http.DefaultClient,
	}
	a.baseUrl, _ = url.Parse(baseUrl)
	a.Headers = make(map[string]string)
	return a
}

func (c *API) NewRequest(method, uri string, opt interface{}, body interface{}) (*http.Request, error) {
	rel, err := addOptions(uri, opt)
	if err != nil {
		return nil, err
	}

	u := c.baseUrl.ResolveReference(rel)
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func (c *API) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		return resp, fmt.Errorf("Server returns status %d", c)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if opt == nil {
		return u, nil
	}

	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		// No query string to add
		return u, nil
	}

	qs, err := query.Values(opt)
	if err != nil {
		return nil, err
	}

	u.RawQuery = qs.Encode()
	return u, nil
}
