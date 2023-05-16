// SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"fmt"
	"testing"
)

func TestClient_ObservationLatestByStationID_Mock(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Station ID
		sid string
		// Altitude
		alt int
		// Latitude
		lat float64
		// Longitude
		lon float64
		// DewPoint
		dp float64
		// DewPoint string
		dps string
	}{
		{"Koeln-Botanischer Garten", "199942", 44, 50.9667, 6.9667, 10.1, "10.1°C"},
		{"Koeln-Stammheim", "H744", 43, 50.9833, 6.9833, 9.7, "9.7°C"},
	}
	c := New(withMockAPI())
	if c == nil {
		t.Errorf("failed to create new Client, got nil")
		return
	}
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			o, err := c.ObservationLatestByStationID(tc.sid)
			if err != nil {
				t.Errorf("ObservationLatestByStationID with station %s failed: %s", tc.sid, err)
				return
			}
			if o.StationID != tc.sid {
				t.Errorf("ObservationLatestByStationID failed, expected station id: %s, got: %s",
					tc.sid, o.StationID)
			}
			if o.Name != tc.n {
				t.Errorf("ObservationLatestByStationID failed, expected name: %s, got: %s",
					tc.n, o.Name)
			}
			if o.Altitude == nil {
				t.Errorf("ObservationLatestByStationID failed, expected altitude but got nil")
			}
			if o.Altitude != nil && *o.Altitude != tc.alt {
				t.Errorf("ObservationLatestByStationID failed, expected altitude: %d, got: %d",
					tc.alt, *o.Altitude)
			}
			if o.Dewpoint() != tc.dps {
				t.Errorf("ObservationLatestByStationID failed, expected dewpoint string: %s, got: %s",
					tc.dps, o.Dewpoint())
			}
		})
	}
}

func TestObservationTemperature_String(t *testing.T) {
	tt := []struct {
		// Original celsius value
		c float64
		// Fahrenheit value
		f float64
	}{
		{-273.15, -459.66999999999996},
		{-50, -58},
		{-40, -40},
		{-30, -22},
		{-20, -4},
		{-10, 14},
		{-9, 15.8},
		{-8, 17.6},
		{-7, 19.4},
		{-6, 21.2},
		{-5, 23},
		{-4, 24.8},
		{-3, 26.6},
		{-2, 28.4},
		{-1, 30.2},
		{0, 32},
		{1, 33.8},
		{2, 35.6},
		{3, 37.4},
		{4, 39.2},
		{5, 41},
		{6, 42.8},
		{7, 44.6},
		{8, 46.4},
		{9, 48.2},
		{10, 50},
		{20, 68},
		{30, 86},
		{40, 104},
		{50, 122},
		{100, 212},
	}
	cf := "%.1f°C"
	ff := "%.1f°F"
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%.2f°C", tc.c), func(t *testing.T) {
			ot := ObservationTemperature{Value: tc.c}
			if ot.Celsius() != tc.c {
				t.Errorf("ObservationTemperature.Celsius failed, expected: %f, got: %f", tc.c,
					ot.Celsius())
			}
			if ot.CelsiusString() != fmt.Sprintf(cf, tc.c) {
				t.Errorf("ObservationTemperature.CelsiusString failed, expected: %s, got: %s",
					fmt.Sprintf(cf, tc.c), ot.CelsiusString())
			}
			if ot.Fahrenheit() != tc.f {
				t.Errorf("ObservationTemperature.Fahrenheit failed, expected: %f, got: %f", tc.f,
					ot.Fahrenheit())
			}
			if ot.FahrenheitString() != fmt.Sprintf(ff, tc.f) {
				t.Errorf("ObservationTemperature.FahrenheitString failed, expected: %s, got: %s",
					fmt.Sprintf(ff, tc.f), ot.FahrenheitString())
			}
		})
	}
}
