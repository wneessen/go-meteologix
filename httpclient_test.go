// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"fmt"
	"testing"
)

// BaseURL is the HTTP Status test base URL
const BaseURL = "https://httpstat.us"

// HTTPStatus is the HTTP Status response type
type HTTPStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func TestNewHTTPClient(t *testing.T) {
	c := New()
	hc := NewHTTPClient(c.config)
	if hc == nil {
		t.Errorf("NewHTTPClient failed, expected HTTPClieht, got nil")
	}
}

func TestHTTPClient_Get(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Status
		s int
		// Expected message
		em string
		// Should fail
		sf bool
	}{
		{"HTTP 200", 200, "OK", false},
		{"HTTP 400", 400, "Bad Request", true},
		{"HTTP 500", 500, "Internal Server Error", true},
	}

	c := New()
	hc := NewHTTPClient(c.config)
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			u := fmt.Sprintf("%s/%d", BaseURL, tc.s)
			r, err := hc.Get(u)
			if err != nil && !tc.sf {
				t.Errorf("HTTPClient Get request failed: %s", err)
				return
			}
			if tc.sf {
				return
			}
			var ro HTTPStatus
			if err := json.Unmarshal(r, &ro); err != nil && !tc.sf {
				t.Errorf("HTTP response unmarshal failed: %s", err)
				return
			}
			if ro.Code != tc.s {
				t.Errorf("HTTPClient Get failed, expected code: %d, got: %d",
					tc.s, ro.Code)
			}
			if ro.Description != tc.em {
				t.Errorf("HTTPClient Get failed, expected message: %s, got: %s",
					tc.em, ro.Description)
			}
		})
	}
}
