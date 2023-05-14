// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"os"
	"testing"
)

func TestClient_StationSearchByLocation(t *testing.T) {
	esid := "199942"
	ak := getAPIKeyFromEnv(t)
	if ak == "" {
		t.Skip("no API_KEY found in environment")
	}
	c := New(WithAPIKey(ak))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	sl, err := c.StationSearchByLocation("Cologne, Germany")
	if err != nil {
		t.Errorf("StationSearchByLocation failed: %s", err)
		return
	}
	if len(sl) < 1 {
		t.Errorf("StationSearchByLocation failed, got no results")
	}
	if sl[0].ID != esid {
		t.Errorf("StationSearchByLocation failed, expected ID: %s, got: %s",
			esid, sl[0].ID)
	}
}

func TestClient_StationSearchByCoordinates(t *testing.T) {
	// RecentlyActive:true, Type:(*string)(nil)
	// Expected stationd data
	p := PrecisionHigh
	es := Station{
		Altitude:       44,
		Distance:       0,
		ID:             "199942",
		Latitude:       50.963,
		Longitude:      6.9698,
		Name:           "Köln-Botanischer Garten",
		Precision:      &p,
		RecentlyActive: true,
	}

	ak := getAPIKeyFromEnv(t)
	if ak == "" {
		t.Skip("no API_KEY found in environment")
	}
	c := New(WithAPIKey(ak))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	sl, err := c.StationSearchByCoordinates(es.Latitude, es.Longitude)
	if err != nil {
		t.Errorf("StationSearchByCoordinates failed: %s", err)
		return
	}
	if len(sl) < 1 {
		t.Errorf("StationSearchByCoordinates failed, got no results")
	}
	rs := sl[0]
	if rs.Altitude != es.Altitude {
		t.Errorf("StationSearchByCoordinates failed, expected altitude: %d, got: %d",
			es.Altitude, rs.Altitude)
	}
	if rs.Distance != es.Distance {
		t.Errorf("StationSearchByCoordinates failed, expected distance: %f, got: %f",
			es.Distance, rs.Distance)
	}
	if rs.ID != es.ID {
		t.Errorf("StationSearchByCoordinates failed, expected id: %s, got: %s",
			es.ID, rs.ID)
	}
	if rs.Latitude != es.Latitude {
		t.Errorf("StationSearchByCoordinates failed, expected latitude: %f, got: %f",
			es.Latitude, rs.Latitude)
	}
	if rs.Longitude != es.Longitude {
		t.Errorf("StationSearchByCoordinates failed, expected longitude: %f, got: %f",
			es.Longitude, rs.Longitude)
	}
	if rs.Name != es.Name {
		t.Errorf("StationSearchByCoordinates failed, expected name: %s, got: %s",
			es.Name, rs.Name)
	}
	if rs.Precision.String() != es.Precision.String() {
		t.Errorf("StationSearchByCoordinates failed, expected precision: %s, got: %s",
			es.Precision, rs.Precision)
	}
	if rs.RecentlyActive != es.RecentlyActive {
		t.Errorf("StationSearchByCoordinates failed, expected recently active: %t, got: %t",
			es.RecentlyActive, rs.RecentlyActive)
	}
}

func TestPrecision_UnmarshalJSON(t *testing.T) {
	type tj struct {
		Precision Precision `json:"precision"`
	}
	tt := []struct {
		// Test name
		n string
		// JSON data
		d []byte
		// Expected precision
		p Precision
		// Should fail
		sf bool
	}{
		{
			"High precision", []byte(`{"precision":"HIGH"}`), PrecisionHigh,
			false,
		},
		{
			"Medium precision", []byte(`{"precision":"MEDIUM"}`), PrecisionMedium,
			false,
		},
		{
			"Low precision", []byte(`{"precision":"LOW"}`), PrecisionLow,
			false,
		},
		{
			"Unknown precision", []byte(`{"precision":"TEST"}`), PrecisionUnknown,
			false,
		},
		{
			"No precision", []byte(`{"precision":null}`), PrecisionUnknown,
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			var p tj
			if err := json.Unmarshal(tc.d, &p); err != nil && !tc.sf {
				t.Errorf("JSON unmarshal failed: %s", err)
				return
			}
			if p.Precision != tc.p {
				t.Errorf("UnmarshalJSON failed, expected: %s, got: %s", tc.p.String(),
					p.Precision.String())
			}
		})
	}
}

func TestPrecision_String(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Precision value
		p Precision
		// Expected string
		es string
	}{
		{"High precision", PrecisionHigh, "HIGH"},
		{"Medium precision", PrecisionMedium, "MEDIUM"},
		{"Low precision", PrecisionLow, "LOW"},
		{"Unknown precision", PrecisionUnknown, "UNKNOWN"},
		{"Undefined precision", 999, "UNKNOWN"},
	}
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			if got := tc.p.String(); got != tc.es {
				t.Errorf("String failed expected %s, got: %s", tc.es, got)
			}
		})
	}
}

func getAPIKeyFromEnv(t *testing.T) string {
	t.Helper()
	return os.Getenv("API_KEY")
}
