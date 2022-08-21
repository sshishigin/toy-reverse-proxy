package simple

import (
	"net/url"
	"time"
)

type simpleServer struct {
	Location  *url.URL
	available bool
	timeout   time.Duration
}

func (s *simpleServer) excludeWithTimeout() {
	s.available = false
	time.Sleep(s.timeout)
	s.available = true
}