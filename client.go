package gohttp

import "net/http"

var DefaultClient = New()

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) do(req *http.Request, transport http.RoundTripper) (*http.Response, error) {
	httpClient := c.httpClient
	if transport != nil {
		httpClient = &http.Client{
			Transport:     transport,
			CheckRedirect: httpClient.CheckRedirect,
			Jar:           httpClient.Jar,
			Timeout:       httpClient.Timeout,
		}
	}
	return httpClient.Do(req)
}

func (c *Client) newRequest(method, url string) *Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return (&Request{}).appendError(err)
	}

	return &Request{
		Request: req,
		client:  c,
	}
}

func (c *Client) Get(url string) *Request {
	return c.newRequest(http.MethodGet, url)
}

func (c *Client) Head(url string) *Request {
	return c.newRequest(http.MethodPost, url)
}

func (c *Client) Post(url string) *Request {
	return c.newRequest(http.MethodPost, url)
}

func (c *Client) Put(url string) *Request {
	return c.newRequest(http.MethodPut, url)
}

func (c *Client) Patch(url string) *Request {
	return c.newRequest(http.MethodPatch, url)
}

func (c *Client) Delete(url string) *Request {
	return c.newRequest(http.MethodDelete, url)
}

func (c *Client) Connect(url string) *Request {
	return c.newRequest(http.MethodConnect, url)
}

func (c *Client) Options(url string) *Request {
	return c.newRequest(http.MethodOptions, url)
}

func (c *Client) Trace(url string) *Request {
	return c.newRequest(http.MethodTrace, url)
}

func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}
