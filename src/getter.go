package src

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	HttpGetter interface {
		HttpGet(uri string) (resultURL string, resp []byte, respErr error)
	}

	httpGetter struct {
		client http.Client
	}
)

const defaultScheme = "http"

var _ HttpGetter = httpGetter{}

func NewHttpGetter(timeoutSec int) HttpGetter {
	return &httpGetter{
		client: http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
	}
}

func (g httpGetter) HttpGet(uri string) (resultURL string, resp []byte, respErr error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", nil, fmt.Errorf("parse URL: %w", err)
	}

	if u.Scheme == "" {
		u.Scheme = defaultScheme
	}
	resultURL = u.String()

	response, err := http.Get(resultURL)
	if err != nil {
		return "", nil, fmt.Errorf("get: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil && respErr == nil {
			respErr = fmt.Errorf("close body: %w", err)
		}
	}()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read body: %w", err)
	}

	return resultURL, responseData, nil
}
