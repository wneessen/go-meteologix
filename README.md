<!--
SPDX-FileCopyrightText: 2023 Winni Neessen <winni@neessen.dev>

SPDX-License-Identifier: CC0-1.0
-->

# go-meteologix
Go bindings to the Meteologix/Kachelmann-Wetter weather API

[![GoDoc](https://godoc.org/github.com/wneessen/go-mail?status.svg)](https://pkg.go.dev/github.com/wneessen/go-meteologix)
[![codecov](https://codecov.io/gh/wneessen/go-meteologix/branch/main/graph/badge.svg?token=W4QI1RMR4L)](https://codecov.io/gh/wneessen/go-meteologix)
[![Go Report Card](https://goreportcard.com/badge/github.com/wneessen/go-meteologix)](https://goreportcard.com/report/github.com/wneessen/go-meteologix)
[![REUSE status](https://api.reuse.software/badge/github.com/wneessen/go-meteologix)](https://api.reuse.software/info/github.com/wneessen/go-meteologix)
<a href="https://ko-fi.com/D1D24V9IX"><img src="https://uploads-ssl.webflow.com/5c14e387dab576fe667689cf/5cbed8a4ae2b88347c06c923_BuyMeACoffee_blue.png" height="20" alt="buy ma a coffee"></a>

## *This package is still WIP*

This Go package provides simple bindings to the 
[Meteologix/Kachelmann-Wetter API](https://api.kachelmannwetter.com/v02/_doc.html#/).
It provides access to "Stations", "Current Weather" and "Forecast". An API key or 
username/password pair is required to access the endpoints. An API key can be configured
in your [account settings](https://accounts.meteologix.com/subscriptions).

go-meteologix follows idiomatic Go style and best practice. It's only dependency is 
the Go Standard Library.

For Geolocation lookups, the package makes use of the 
[OpenStreetMap Nominatim API](https://nominatim.org/). This requires no API key.

## Usage

The library is fully documented using the execellent GoDoc functionality. Check out 
the [full reference on pkg.go.dev](https://pkg.go.dev/github.com/wneessen/go-hibp) for 
details.

## Examples

### GeoLocation lookup

This program uses the OSM Nominatim API to lookup the GeoLocation data for *Berlin, Germany*.
On success it will return the `Latitude` and `Longitude` fields.
```go
package main

import (
	"fmt"
	"os"

	"github.com/wneessen/go-meteologix"
)

func main() {
	c := meteologix.New()
	gl, err := c.GetGeoLocationByName("Berlin, Germany")
	if err != nil {
		fmt.Println("GeoLocation lookup failed", err)
		os.Exit(1)
	}
	fmt.Printf("GeoLocation - Latitude: %f, Longitude: %f\n", gl.Latitude,
		gl.Longitude)
}
```
