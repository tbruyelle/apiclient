package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *API) NewRequest(method, uri string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(uri)
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
