package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type registered struct {
	sync.RWMutex
	uuid map[string]*Session // registered user[uuid]session
	id   map[string]*Session // registered user[user_id]session
}

func (r *registered) actualizeDate(w http.ResponseWriter, id string) object.Status {
	if r.id[id].expired() {
		// r.UpdateOrBanRS()
	}
	return nil
}

// CreateRS creates session for registered user
func (r *registered) CreateRS(w http.ResponseWriter, userID string) object.Status {
	// generate uuid
	uuid, sts := genUUID()
	if sts != nil {
		return sts
	}
	// add session
	r.makeRS(uuid, userID)
	// create cookie
	// create cookie
	src := map[string]string{
		config.CookieSession: uuid,
		config.CookieUserID:  userID,
	}
	object.NewCookie(w, src).Set(true)
	return nil
}

func (r *registered) makeRS(uuid, id string) {
	now := time.Now()
	ses := &Session{
		uuid:   uuid,
		id:     id,
		expire: now,
		requests: requests{
			req:   [config.CheckAfter]time.Time{now},
			count: 1,
		},
	}
	r.Lock()
	defer r.Unlock()
	r.uuid[uuid] = ses
	r.id[id] = ses
}

func (r *registered) UpdateOrBanRS(w http.ResponseWriter, userID string) (object.Status, bool) {
	// generate uuid
	uuid, sts := genUUID()
	if sts != nil {
		return sts, false
	}
	// add session
	r.updateRS(uuid, userID)
	// create cookie
	// create cookie
	src := map[string]string{
		config.CookieSession: uuid,
		config.CookieUserID:  userID,
	}
	object.NewCookie(w, src).Set(true)
	return nil, false
}

func (r *registered) updateRS(uuid, id string) {
	r.Lock()
	defer r.Unlock()
	s := r.id[id]
	s.uuid = uuid
	s.expire = time.Now().Add(config.SessionExpire)
	s.count++
	// if s.isDDos() {
	// }

	r.uuid[uuid] = s
}

func (r *registered) match(ck *object.CookieInfo) bool {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.uuid[ck.UUID]; ok && v.id == ck.UserStr {
		return true
	}
	return false
}

func (r *registered) delete(key string) {
	r.RLock()
	defer r.RUnlock()
	if v, ok := r.uuid[key]; ok {
		id := v.id
		r.deleteUUID(key)
		r.deleteID(id)
	} else if v, ok := r.id[key]; ok {
		uuid := v.uuid
		r.deleteID(key)
		r.deleteUUID(uuid)
	}
}

func (r *registered) deleteUUID(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.uuid, key)
}

func (r *registered) deleteID(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.id, key)
}
