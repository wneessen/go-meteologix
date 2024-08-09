// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// HTTPClientTimeout is the default timeout value for the HTTPClient
	HTTPClientTimeout = time.Second * 10
	// MIMETypeJSON is a string constant for application/json MIME type
	MIMETypeJSON = "application/json"
)

// ErrNonJSONResponse is returned when a HTTPClient request did not return the expected
// application/json content type
var ErrNonJSONResponse = errors.New("HTTP response is of non-JSON content type")

// HTTPClient is a type wrapper for the Go stdlib http.Client and the Config
type HTTPClient struct {
	*Config
	*http.Client
}

// APIError wraps the error interface for the API
type APIError struct {
	Code    int    `json:"status"`
	Details string `json:"detail"`
	Message string `json:"message"`
	Title   string `json:"title"`
	Type    string `json:"type"`
}

// NewHTTPClient returns a new HTTP client
func NewHTTPClient(config *Config) *HTTPClient {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	httpTransport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{
		Timeout:   HTTPClientTimeout,
		Transport: httpTransport,
	}
	return &HTTPClient{config, httpClient}
}

// Get performs a HTTP GET request for the given URL with the default HTTP timeout
func (hc *HTTPClient) Get(url string) ([]byte, error) {
	return hc.GetWithTimeout(url, HTTPClientTimeout)
}

// GetWithTimeout performs a HTTP GET request for the given URL and sets a timeout context
// with the given timeout duration
func (hc *HTTPClient) GetWithTimeout(url string, timeout time.Duration) ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", hc.userAgent)
	request.Header.Set("Content-Type", MIMETypeJSON)
	request.Header.Set("Accept", MIMETypeJSON)
	request.Header.Set("Accept-Language", hc.acceptLang)

	// User authentication (only required for Meteologix API calls)
	if strings.HasPrefix(url, APIBaseURL) {
		hc.setAuthentication(request)
	}

	response, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("nil response received")
	}
	defer func(body io.ReadCloser) {
		if err = body.Close(); err != nil {
			log.Printf("failed to close HTTP request body: %s", err)
		}
	}(response.Body)

	if !strings.HasPrefix(response.Header.Get("Content-Type"), MIMETypeJSON) {
		return nil, ErrNonJSONResponse
	}
	if response.StatusCode >= http.StatusBadRequest {
		apiError := new(APIError)
		if err = json.NewDecoder(response.Body).Decode(apiError); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error JSON: %w", err)
		}
		if apiError.Code < 1 {
			apiError.Code = response.StatusCode
		}
		if apiError.Details == "" {
			apiError.Details = response.Status
		}
		return nil, *apiError
	}

	buffer := &bytes.Buffer{}
	bufferWriter := bufio.NewWriter(buffer)
	_, err = io.Copy(bufferWriter, response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy HTTP response body to buffer: %w", err)
	}
	if err = bufferWriter.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %w", err)
	}
	return buffer.Bytes(), nil
}

// setAuthentication sets the corresponding user authentication header. If an API Key is set, this
// will be preferred, alternatively a username/authPass combination for HTTP Basic auth can
// be used
func (hc *HTTPClient) setAuthentication(httpRequest *http.Request) {
	if hc.apiKey != "" {
		httpRequest.Header.Set("X-API-Key", hc.Config.apiKey)
		return
	}
	if hc.bearerToken != "" {
		httpRequest.Header.Set("Authorization", "Bearer"+hc.bearerToken)
		return
	}
	if hc.authUser != "" && hc.authPass != "" {
		httpRequest.SetBasicAuth(url.QueryEscape(hc.authUser), url.QueryEscape(hc.authPass))
	}
}

// Error satisfies the error interface for the APIError type
func (e APIError) Error() string {
	var errorMsg strings.Builder
	errorMsg.WriteString("API request failed with status HTTP ")
	errorMsg.WriteString(fmt.Sprintf("%d: ", e.Code))
	if e.Details != "" {
		errorMsg.WriteString(e.Details)
	}
	if e.Message != "" {
		errorMsg.WriteString(" (Optional message: ")
		errorMsg.WriteString(e.Message)
		errorMsg.WriteString(")")
	}
	return errorMsg.String()
}
