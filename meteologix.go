// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"runtime"
)

const (
	// APIBaseURL represents the base URL for the Meteologix API.
	//
	// We currently support v02 of the API.
	APIBaseURL = "https://api.kachelmannwetter.com/v02"
	// APIMockURL represents the mocked API URL for testing purposes
	APIMockURL = "https://go-meteologix-mock.neessen.dev/v02"
	// DefaultAcceptLang is the default language set for API requests
	DefaultAcceptLang = "en"
)

// DefaultUserAgent is the default User-Agent presented by the HTTPClient
var DefaultUserAgent = fmt.Sprintf("go-meteologix/v%s (%s; %s; "+
	"+https://github.com/wneessen/go-meteologix)", VERSION, runtime.GOOS,
	runtime.Version())

// Client represents the Meteologix API Client
type Client struct {
	// config represents the Config for the Client
	config *Config
	// httpClient references the HTTPClient of the Server
	httpClient *HTTPClient
}

// Config represents the Client configuration settings
type Config struct {
	// apiKey holds the (optional) API key for the API user authentication
	apiKey string
	// apiURL holds the base URL for the API. This is configurable so we
	// can test against our mock API.
	apiURL string
	// acceptLang hold the (optional) accept-language tag
	acceptLang string
	// authPass holds the (optional) passowrd for the API user authentication
	authPass string
	// authUser holds the (optional) username for the API user authentication
	authUser string
	// userAgent represents an alternative User-Agent HTTP header string
	userAgent string
}

// Option represents a function that is used for setting/overriding Config options
type Option func(*Config)

// New returns a new Meteologix API Client
func New(o ...Option) *Client {
	co := &Config{}
	co.apiURL = APIBaseURL
	co.acceptLang = DefaultAcceptLang
	co.userAgent = DefaultUserAgent

	// Set/override Config options
	for _, opt := range o {
		if opt == nil {
			continue
		}
		opt(co)
	}

	return &Client{
		config:     co,
		httpClient: NewHTTPClient(co),
	}
}

// WithAcceptLanguage sets the HTTP Accept-Lanauge header of the HTTP client
//
// The provided string needs to conform the HTTP Accept-Language header format
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func WithAcceptLanguage(l string) Option {
	if l == "" {
		return nil
	}
	return func(co *Config) {
		co.acceptLang = l
	}
}

// WithAPIKey sets the API Key for user authentication of the HTTP client
func WithAPIKey(k string) Option {
	if k == "" {
		return nil
	}
	return func(co *Config) {
		co.apiKey = k
	}
}

// WithPassword sets the HTTP Basic auth authPass for the HTTP client
func WithPassword(p string) Option {
	if p == "" {
		return nil
	}
	return func(co *Config) {
		co.authPass = p
	}
}

// WithUserAgent sets a custom user agent string for the HTTP client
func WithUserAgent(a string) Option {
	if a == "" {
		return nil
	}
	return func(co *Config) {
		co.userAgent = a
	}
}

// WithUsername sets the HTTP Basic auth username for the HTTP client
func WithUsername(u string) Option {
	if u == "" {
		return nil
	}
	return func(co *Config) {
		co.authUser = u
	}
}

// withMockAPI sets the API URL to our mock API for testing
func withMockAPI() Option {
	return func(co *Config) {
		co.apiURL = APIMockURL
	}
}
