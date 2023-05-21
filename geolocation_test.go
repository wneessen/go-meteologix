// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT

package meteologix

import (
	"errors"
	"strings"
	"testing"
)

func TestClient_GetGeoLocationByCityName(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Query string
		q string
		// Expected name
		en []string
		// Expected latitude
		ela float64
		// Expected longitude
		elo float64
		// Should fail
		sf bool
	}{
		{
			"Cologne DE", "Cologne, Germany",
			[]string{"Cologne", "North Rhine-Westphalia", "Germany"},
			50.938361,
			6.959974, false,
		},
		{
			"Neermoor DE", "Neermoor, Germany",
			[]string{"Neermoor", "Moormerland", "Landkreis Leer", "Germany"},
			53.3067155, 7.4418682, false,
		},
		{
			"Chicago US", "Chicago",
			[]string{"Chicago", "Cook County", "Illinois", "United States"},
			41.8755616, -87.6244212, false,
		},
		{"Nonexisting", "Nonexisting City", []string{}, 0, 0, true},
	}

	c := New()
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			l, err := c.GetGeoLocationByName(tc.q)
			if err != nil && !tc.sf {
				t.Errorf("GetGeoLocationByName failed: %s", err)
			}
			if tc.sf {
				return
			}
			fn := 0
			for i := range tc.en {
				if strings.Contains(l.Name, tc.en[i]) {
					fn++
				}
			}
			if len(tc.en) != fn {
				t.Errorf("GetGeoLocationByName failed, expected %d matching name parts, got: %d",
					len(tc.en), fn)
			}
			if l.Latitude != tc.ela {
				t.Errorf("GetGeoLocationByName failed, expected latitude: %f, got: %f", tc.ela, l.Latitude)
			}
			if l.Longitude != tc.elo {
				t.Errorf("GetGeoLocationByName failed, expected longitude: %f, got: %f", tc.elo, l.Longitude)
			}
		})
	}
}

func TestClient_GetGeoLocationByCityName_CityNotFoundErr(t *testing.T) {
	c := New()
	_, err := c.GetGeoLocationByName("Nonexisting City")
	if err == nil {
		t.Errorf("GetGeoLocationByName with nonexisting city name was expected to fail, but didn't")
		return
	}
	if !errors.Is(err, ErrCityNotFound) {
		t.Errorf("GetGeoLocationByName was supposed to fail with ErrCityNotFound error, but didn't")
	}
}
