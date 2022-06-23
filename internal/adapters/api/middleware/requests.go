package middleware

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
)

// save requests from one user_id/ip
type requests struct {
	req   [config.CheckAfter]time.Time // for request and time when came
	count int                          // count added requests
}

func newRequests() *requests {
	return &requests{
		req:   [config.CheckAfter]time.Time{time.Now()},
		count: 1,
	}
}

// isDDos checks how many requests came from one user_id/ip (parent structure)
func (r *requests) isDDos() (bool, float64) {
	if r.count < config.CheckAfter {
		r.req[r.count-1] = time.Now()
		return false, 0
	}
	// if saved requests == "check_after"
	// time to check requsts interval for ddos
	firstReq := r.req[0]
	lastReq := r.req[config.CheckAfter-1]
	interval := lastReq.Sub(firstReq).Seconds()
	return interval <= config.CheckTimeInterval.Seconds(), interval
}