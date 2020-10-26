package gohttp

import "net/http"

func Get(url string) *Request {
	return DefaultClient.newRequest(http.MethodGet, url)
}

func Head(url string) *Request {
	return DefaultClient.newRequest(http.MethodPost, url)
}

func Post(url string) *Request {
	return DefaultClient.newRequest(http.MethodPost, url)
}

func Put(url string) *Request {
	return DefaultClient.newRequest(http.MethodPut, url)
}

func Patch(url string) *Request {
	return DefaultClient.newRequest(http.MethodPatch, url)
}

func Delete(url string) *Request {
	return DefaultClient.newRequest(http.MethodDelete, url)
}

func Connect(url string) *Request {
	return DefaultClient.newRequest(http.MethodConnect, url)
}

func Options(url string) *Request {
	return DefaultClient.newRequest(http.MethodOptions, url)
}

func Trace(url string) *Request {
	return DefaultClient.newRequest(http.MethodTrace, url)
}
