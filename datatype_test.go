// SPDX-FileCopyrightText: 2023 Winni Neessen <wn@neessen.dev>
//
// SPDX-License-Identifier: MIT
package meteologix

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAPIDate_UnmarshalJSON(t *testing.T) {
	type testType struct {
		Date APIDate `json:"date"`
	}
	f := func(jsonData []byte, expected string, shouldFail bool) {
		t.Helper()

		var date testType
		if err := json.Unmarshal(jsonData, &date); err != nil && !shouldFail {
			t.Errorf("APIDate_UnmarshalJSON failed: %s", err)
			return
		}
		if !strings.EqualFold(date.Date.Format(DateFormat), expected) && !shouldFail {
			t.Errorf("APIDate_UnmarshalJSON failed, expected: %s, but got: %s",
				expected, date.Date.Format(DateFormat))
		}
	}

	f([]byte(`{"date":"2023-05-28"}`), "2023-05-28", false)
	f([]byte(`{"date":"2023-05-32"}`), "2023-05-32", true)
	f([]byte(`{"date":null}`), "", true)
}
