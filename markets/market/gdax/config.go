package gdax

import (
	"net/url"
	"time"
)

var (
	DefaultConfig = Config{
		Endpoint: EndpointConfig{
			URL: &url.URL{
				Scheme: "wss",
				Host:   "ws-feed.gdax.com",
				Path:   "/",
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
