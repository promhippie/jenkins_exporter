package exporter

import (
	"crypto/tls"
	"net/http"

	"github.com/jackspirou/syscerts"
)

// simpleClient initializes a simple HTTP client.
func simpleClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs: syscerts.SystemRootsPool(),
			},
		},
	}
}

// colorToGauge maps the job color to a float.
func colorToGauge(color string) float64 {
	switch color {
	case "blue":
		return 0.0
	case "blue_anime":
		return 0.5
	case "red":
		return 1.0
	case "red_anime":
		return 1.5
	case "yellow":
		return 1.0
	case "yellow_anime":
		return 1.5
	case "notbuilt":
		return 2.0
	case "notbuilt_anime":
		return 2.5
	case "disabled":
		return 3.0
	case "disabled_anime":
		return 3.5
	case "aborted":
		return 4.0
	case "aborted_anime":
		return 4.5
	case "grey":
		return 5.0
	case "grey_anime":
		return 5.5
	}

	return -1.0
}
