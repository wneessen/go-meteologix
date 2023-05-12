// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

// APIBaseURL represents the base URL for the Meteologix API
// We currently support v02 of the API
const APIBaseURL = "https://api.kachelmannwetter.com/v02"

const (
	// DefaultAcceptLang is the default language set for API requests
	DefaultAcceptLang = "en"
	// DefaultUserAgent is the default User-Agent presented by the HTTPClient
	DefaultUserAgent = "go-meteologix v" + VERSION
)

// Client represents the Meteologix API Client
type Client struct {
	// co represents the Config for the Client
	co *Config
	// hc references the HTTPClient of the Server
	hc *HTTPClient
}

// Config represents the Client configuration settings
type Config struct {
	// ak holds the (optional) API key for the API user authentication
	ak string
	// al hold the (optional) accept-language tag
	al string
	// pw holds the (optional) passowrd for the API user authentication
	pw string
	// ua represents an alternative User-Agent HTTP header string
	ua string
	// un holds the (optional) username for the API user authentication
	un string
}

// Option represents a function that is used for setting/overriding Config options
type Option func(*Config)

// New returns a new Meteologix API Client
func New(o ...Option) *Client {
	co := &Config{}
	co.al = DefaultAcceptLang
	co.ua = DefaultUserAgent

	// Set/override Config options
	for _, opt := range o {
		if opt == nil {
			continue
		}
		opt(co)
	}

	return &Client{
		co: co,
		hc: NewHTTPClient(co),
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
		co.al = l
	}
}

// WithAPIKey sets the API Key for user authentication of the HTTP client
func WithAPIKey(k string) Option {
	if k == "" {
		return nil
	}
	return func(co *Config) {
		co.ak = k
	}
}

// WithPassword sets the HTTP Basic auth password for the HTTP client
func WithPassword(p string) Option {
	if p == "" {
		return nil
	}
	return func(co *Config) {
		co.pw = p
	}
}

// WithUserAgent sets a custom user agent string for the HTTP client
func WithUserAgent(a string) Option {
	if a == "" {
		return nil
	}
	return func(co *Config) {
		co.ua = a
	}
}

// WithUsername sets the HTTP Basic auth username for the HTTP client
func WithUsername(u string) Option {
	if u == "" {
		return nil
	}
	return func(co *Config) {
		co.un = u
	}
}
