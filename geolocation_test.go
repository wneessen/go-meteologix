package meteologix

import (
	"errors"
	"testing"
)

func TestClient_GetGeoLocationByCityName(t *testing.T) {
	tt := []struct {
		// Test name
		n string
		// Query string
		q string
		// Expected name
		en string
		// Expected latitude
		ela float64
		// Expected longitude
		elo float64
		// Should fail
		sf bool
	}{
		{
			"Cologne, DE", "Cologne, Germany", "Cologne, North Rhine-Westphalia, Germany", 50.938361,
			6.959974, false,
		},
		{
			"Neermoor, DE", "Neermoor, Germany",
			"Neermoor, Moormerland, Landkreis Leer, Lower Saxony, 26802, Germany", 53.3067155,
			7.4418682, false,
		},
		{
			"Chicago, US", "Chicago", "Chicago, Cook County, Illinois, United States", 41.8755616,
			-87.6244212, false,
		},
		{"Nonexisting", "Nonexisting City", "", 0, 0, true},
	}

	c := New()
	for _, tc := range tt {
		t.Run(tc.n, func(t *testing.T) {
			l, err := c.GetGeoLocationByCityName(tc.q)
			if err != nil && !tc.sf {
				t.Errorf("GetGeoLocationByCityName failed: %s", err)
			}
			if tc.sf {
				return
			}
			if l.Name != tc.en {
				t.Errorf("GetGeoLocationByCityName failed, expected name: %s, got: %s", tc.en, l.Name)
			}
			if l.Latitude != tc.ela {
				t.Errorf("GetGeoLocationByCityName failed, expected latitude: %f, got: %f", tc.ela, l.Latitude)
			}
			if l.Longitude != tc.elo {
				t.Errorf("GetGeoLocationByCityName failed, expected longitude: %f, got: %f", tc.elo, l.Longitude)
			}
		})
	}
}

func TestClient_GetGeoLocationByCityName_CityNotFoundErr(t *testing.T) {
	c := New()
	_, err := c.GetGeoLocationByCityName("Nonexisting City")
	if err == nil {
		t.Errorf("GetGeoLocationByCityName with nonexisting city name was expected to fail, but didn't")
		return
	}
	if !errors.Is(err, ErrCityNotFound) {
		t.Errorf("GetGeoLocationByCityName was supposed to fail with ErrCityNotFound error, but didn't")
	}
}
