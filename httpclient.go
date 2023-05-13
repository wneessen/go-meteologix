// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// HTTPClientTimeout is the default timeout value for the HTTPClient
	HTTPClientTimeout = time.Second * 10
	// MIMETypeJSON is a string constant for application/json MIME type
	MIMETypeJSON = "application/json"
)

// ErrNonJSONResponse is returned when a HTTPClient request did not return the
// expected application/json content type
var ErrNonJSONResponse = errors.New("HTTP response is of non-JSON content type")

// HTTPClient is a type wrapper for the Go stdlib http.Client and the Config
type HTTPClient struct {
	*Config
	*http.Client
}

// NewHTTPClient returns a new HTTP client
func NewHTTPClient(c *Config) *HTTPClient {
	tc := &tls.Config{
		MaxVersion: tls.VersionTLS12,
		MinVersion: tls.VersionTLS12,
	}
	ht := http.Transport{TLSClientConfig: tc}
	hc := &http.Client{Transport: &ht}
	return &HTTPClient{c, hc}
}

// Get performs a HTTP GET request for the given URL with the default HTTP timeout
func (hc *HTTPClient) Get(u string) ([]byte, error) {
	return hc.GetWithTimeout(u, HTTPClientTimeout)
}

// GetWithTimeout performs a HTTP GET request for the given URL and sets a timeout
// context with the given timeout duration
func (hc *HTTPClient) GetWithTimeout(u string, t time.Duration) ([]byte, error) {
	ctx, cfn := context.WithTimeout(context.Background(), t)
	defer cfn()
	hr, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	hr.Header.Set("User-Agent", hc.userAgent)
	hr.Header.Set("Content-Type", MIMETypeJSON)
	hr.Header.Set("Accept", MIMETypeJSON)
	hr.Header.Set("Accept-Language", hc.acceptLang)

	// User authentication (only required for Meteologix API calls)
	if strings.HasPrefix(u, APIBaseURL) {
		hc.setAuthHeader(hr)
	}

	sr, err := hc.Do(hr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := sr.Body.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "failed to close HTTP request body", err)
		}
	}()

	if !strings.HasPrefix(sr.Header.Get("Content-Type"), MIMETypeJSON) {
		return nil, ErrNonJSONResponse
	}
	buf := &bytes.Buffer{}
	bw := bufio.NewWriter(buf)
	_, err = io.Copy(bw, sr.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy HTTP response body to buffer: %w", err)
	}
	return buf.Bytes(), nil
}

// setAuthHeader sets the corresponding user authentication header. If an API Key is set, this
// will be preferred, alternatively a username/authPass combination for HTTP Basic auth can
// be used
func (hc *HTTPClient) setAuthHeader(hr *http.Request) {
	if hc.apiKey != "" {
		hr.Header.Set("X-API-Key", hc.Config.apiKey)
		return
	}
	if hc.authUser != "" && hc.authPass != "" {
		hr.SetBasicAuth(hc.authUser, hc.authPass)
	}
}
