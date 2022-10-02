package weighted

import (
	"net/url"
	"time"
)

type weightedServer struct {
	Location  *url.URL
	available bool
	timeout   time.Duration
	weight    int
}

func (s *weightedServer) excludeWithTimeout() {
	s.available = false
	time.Sleep(s.timeout)
	s.available = true
}
