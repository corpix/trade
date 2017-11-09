package bitfinex

import (
	"net/url"
	"time"
)

var (
	DefaultConfig = Config{
		Endpoint: EndpointConfig{
			URL: &url.URL{
				Scheme: "wss",
				Host:   "api.bitfinex.com",
				Path:   "/ws/2",
			},
			Timeout: 5 * time.Second,
		},
	}
)

type Config struct {
	Token    string
	Endpoint EndpointConfig
}

type EndpointConfig struct {
	URL     *url.URL
	Timeout time.Duration
}
