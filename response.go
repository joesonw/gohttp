package gohttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	read  bool
	bytes []byte
	err   error
}

func (res *Response) OK() bool {
	return res.StatusCode < 300
}

func (res *Response) HTTPError() error {
	if res.OK() {
		return nil
	}

	bytes, err := res.Bytes()
	if err != nil {
		return err
	}

	return fmt.Errorf("%s: %s", res.Status, string(bytes))
}

func (res *Response) Error() error {
	return res.err
}

func (res *Response) Bytes() ([]byte, error) {
	if res.err != nil {
		return nil, res.err
	}

	if res.read {
		return res.bytes, nil
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	res.read = true
	res.bytes = bytes
	return bytes, err
}

func (res *Response) ToJSON(in interface{}) error {
	bytes, err := res.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, in)
}

func (res *Response) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	if err := res.ToJSON(&m); err != nil {
		return nil, err
	}
	return m, nil
}
