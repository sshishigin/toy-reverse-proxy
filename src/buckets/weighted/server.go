package weighted

import (
	"net/url"
	"time"
)

type Server struct {
	Location  *url.URL
	available bool
	timeout   time.Duration
	weight    int
}

func (s *Server) ExcludeWithTimeout() {
	s.available = false
	time.Sleep(s.timeout)
	s.available = true
}
