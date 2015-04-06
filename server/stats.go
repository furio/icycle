package server

import (
    "fmt"
    "time"
)

type Stats struct {
    requests uint64
    requestsTime uint64
}

func NewStats() (*Stats) {
    s := &Stats{}
    s.requests = uint64(0)
    s.requestsTime = uint64(0)

    return s
}

func (s *Stats) RecordRequest(startTime time.Time) {
    s.requests += 1
    s.requestsTime += uint64(time.Since(startTime).Nanoseconds())
}

func (s *Stats) TotalStats() string {
    return fmt.Sprintf("Requests: %d, Time: %d ns, Avg: %d ns/req", s.requests, s.requestsTime, s.requestsTime/s.requests)
}