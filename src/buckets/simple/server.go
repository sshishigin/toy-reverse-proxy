package simple

import (
	"log"
	"net/url"
	"time"
)

type Server struct {
	Location  *url.URL
	available bool
	timeout   time.Duration
	maxFails  int
	fails     int
}

func (s *Server) ExcludeWithTimeout() {
	if s.available {
		s.available = false
		log.Printf("excluding %s from bucket for %s", s.Location, s.timeout)
		time.Sleep(s.timeout)
		s.available = true
		s.fails = 0
		log.Printf("%s is back", s.Location)
	}
}
