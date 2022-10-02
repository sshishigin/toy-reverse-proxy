package config

import (
	"net/url"
	"time"
)

type BucketConfig struct {
	Servers  []ServerConfig `json:"servers"`
	MaxFails int            `json:"maxFails"`
	Timeout  time.Duration  `json:"timeout"`
}

type ServerConfig struct {
	Location *url.URL `json:"location"`
	Weight   int      `json:"weight"`
}
