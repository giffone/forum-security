package middleware

import (
	"log"
	"sync"
	"time"

	"github.com/giffone/forum-security/internal/config"
)

type unRegistered struct {
	sync.RWMutex
	ip map[string]*requests // unregistered user[ip]session
}

func (ur *unRegistered) make(key string) bool {
	// check exist or no
	ur.RLock()
	_, ok := ur.ip[key]
	ur.RUnlock()
	if !ok {
		ur.add(key) // create new map
		return true
	}
	ur.RLock()
	ddos, interval := ur.ip[key].isDDos()
	ur.RUnlock()
	if ddos {
		log.Printf("ddos with interval: %.4f with limit %.4f\n", interval, config.CheckTimeInterval.Seconds())
		return !ddos
	}
	ur.clear(key)
	return true
}

func (ur *unRegistered) add(key string) {
	req := &requests{
		req:   [config.CheckAfter]time.Time{time.Now()},
		count: 1,
	}
	ur.Lock()
	defer ur.Unlock()
	ur.ip[key] = req
}

func (ur *unRegistered) delete(key string) {
	ur.Lock()
	defer ur.Unlock()
	delete(ur.ip, key)
}

func (ur *unRegistered) clear(key string) {
	ur.Lock()
	defer ur.Unlock()
	ur.ip[key] = newRequests()
}
