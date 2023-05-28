// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT
package meteologix

import (
	"encoding/json"
	"testing"
)

func TestAPIDate_UnmarshalJSON(t *testing.T) {
	type x struct {
		Date APIDate `json:"date"`
	}
	okd := []byte(`{"date":"2023-05-28"}`)
	nokd := []byte(`{"date":"2023-05-32"}`)
	null := []byte(`{"date":null}`)
	var d x
	if err := json.Unmarshal(okd, &d); err != nil {
		t.Errorf("APIDate_UnmarshalJSON failed: %s", err)
	}
	if d.Date.Format(DateFormat) != "2023-05-28" {
		t.Errorf("APIDate_UnmarshalJSON failed, expected: %s, but got: %s",
			"2023-05-28", d.Date.String())
	}
	if err := json.Unmarshal(nokd, &d); err == nil {
		t.Errorf("APIDate_UnmarshalJSON was supposed to fail, but didn't")
	}
	d = x{}
	if err := json.Unmarshal(null, &d); err != nil {
		t.Errorf("APIDate_UnmarshalJSON failed: %s", err)
	}
	if !d.Date.IsZero() {
		t.Errorf("APIDate_UnmarshalJSON with null was supposed to be empty, but got: %s",
			d.Date.String())
	}
}
