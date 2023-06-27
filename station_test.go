// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"encoding/json"
	"errors"
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

func TestClient_StationSearchByLocation_Fail(t *testing.T) {
	c := New(WithUsername("invalid"), WithPassword("invalid"))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	_, err := c.StationSearchByLocation("Cologne, Germany")
	if err == nil {
		t.Errorf("StationSearchByLocation was supposed to fail but didn't")
	}
	if err != nil && !errors.As(err, &APIError{}) {
		t.Errorf("StationSearchByLocation was supposed to throw a APIError but didn't: %s",
			err)
	}
	c = New(WithAPIKey("invalid"))
	_, err = c.StationSearchByLocation("Cologne, Germany")
	if err == nil {
		t.Errorf("StationSearchByLocation was supposed to fail but didn't")
		return
	}
	if err != nil && !errors.As(err, &APIError{}) {
		t.Errorf("StationSearchByLocation was supposed to throw a APIError but didn't: %s",
			err)
	}
}

func TestClient_StationSearchByLocationWithRadius_Fail(t *testing.T) {
	ak := getAPIKeyFromEnv(t)
	if ak == "" {
		t.Skip("no API_KEY found in environment, skipping test")
	}
	c := New(WithAPIKey(ak))
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	_, err := c.StationSearchByLocationWithinRadius("Cologne, Germany", 0)
	if err == nil {
		t.Errorf("StationSearchByLocationWithRadius was supposed to fail but didn't")
	}
	if !errors.Is(err, ErrRadiusTooSmall) {
		t.Errorf("StationSearchByLocationWithRadius was supposed to return ErrRadiusTooSmall, got: %s", err)
	}
	_, err = c.StationSearchByLocationWithinRadius("Cologne, Germany", 1000)
	if err == nil {
		t.Errorf("StationSearchByLocationWithRadius was supposed to fail but didn't")
	}
}

func TestClient_StationSearchByCoordinates_Mock(t *testing.T) {
	// Expected station data from mock API
	p := PrecisionHigh
	ty := "STATDEU6"
	es := Station{
		Altitude:       822,
		Distance:       12.6,
		ID:             "106350",
		Latitude:       50.221,
		Longitude:      8.4469,
		Name:           "Feldberg/Taunus",
		Precision:      &p,
		RecentlyActive: true,
		Type:           &ty,
	}

	c := New(withMockAPI())
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
	if *rs.Type != *es.Type {
		t.Errorf("StationSearchByCoordinates failed, expected type: %s, got: %s",
			*es.Type, *rs.Type)
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
			"Super high precision", []byte(`{"precision":"SUPER_HIGH"}`), PrecisionSuperHigh,
			false,
		},
		{
			"High precision", []byte(`{"precision":"HIGH"}`), PrecisionHigh,
			false,
		},
		{
			"Standard precision", []byte(`{"precision":"STANDARD"}`), PrecisionStandard,
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
		{"Super high precision", PrecisionSuperHigh, "SUPER_HIGH"},
		{"High precision", PrecisionHigh, "HIGH"},
		{"Standard precision", PrecisionStandard, "STANDARD"},
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
