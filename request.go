package gohttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.uber.org/multierr"
)

type Request struct {
	*http.Request
	client    *Client
	errs      []error
	query     url.Values
	tlsConfig *tls.Config
}

func (req *Request) Error() error {
	return multierr.Combine(req.errs...)
}

func (req *Request) appendError(err error) *Request {
	req.errs = append(req.errs, err)
	return req
}

func (req *Request) ensureTLSConfig() {
	if req.tlsConfig == nil {
		req.tlsConfig = &tls.Config{}
	}
}

func (req *Request) Do() *Response {
	if err := req.Error(); err != nil {
		return &Response{err: err}
	}

	if len(req.query) > 0 {
		q := req.Request.URL.Query()
		for k, a := range req.query {
			for _, v := range a {
				q[k] = append(q[k], v)
			}
		}
		req.Request.URL.RawQuery = q.Encode()
	}

	var transport http.RoundTripper
	if req.tlsConfig != nil {
		transport = &http.Transport{
			TLSClientConfig: req.tlsConfig,
		}
	}

	res, err := req.client.do(req.Request, transport)
	if err != nil {
		return &Response{err: err}
	}

	return &Response{
		Response: res,
	}
}

func (req *Request) WithContext(ctx context.Context) *Request {
	req.Request = req.Request.WithContext(ctx)
	return req
}

func (req *Request) Query(key, value string) *Request {
	req.query.Set(key, value)
	return req
}

func (req *Request) Set(key, value string) *Request {
	req.Header.Set(key, value)
	return req
}

func (req *Request) JSON(in interface{}) *Request {
	b, err := json.Marshal(in)
	if err != nil {
		return req.appendError(err)
	}
	req.Set("Content-Type", "application/json")
	req.setBody(b)
	return req
}

func (req *Request) Form(values url.Values) *Request {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.setBody([]byte(values.Encode()))
	return req
}

func (req *Request) setBody(b []byte) {
	reader := bytes.NewReader(b)
	req.Request.Body = ioutil.NopCloser(reader)
	req.Request.ContentLength = int64(len(b))
	req.Request.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(reader), nil
	}
}

func (req *Request) Clone(ctx context.Context) *Request {
	httpRequest := req.Request.Clone(ctx)
	query := url.Values{}
	for k, a := range req.query {
		for _, v := range a {
			query[k] = append(query[k], v)
		}
	}

	newReq := &Request{
		Request: httpRequest,
		client:  req.client,
		errs:    req.errs[:],
		query:   query,
	}

	return newReq
}

func (req *Request) WithTLSConfig(tlsConfig *tls.Config) *Request {
	req.tlsConfig = tlsConfig
	return req
}

func (req *Request) Insecure(insecure bool) *Request {
	req.ensureTLSConfig()
	req.tlsConfig.InsecureSkipVerify = insecure
	return req
}

func (req *Request) WithCertFile(file string) *Request {
	caCert, err := ioutil.ReadFile(file)
	if err != nil {
		return req.appendError(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	req.tlsConfig.RootCAs = caCertPool
	return req
}
