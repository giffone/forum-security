package middleware

import (
	"net/http"
	"time"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/adapters/api/middleware/bann"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	uuid "github.com/nu7hatch/gouuid"
)

type Ses struct {
	reg    registered   // registered users and its sessions and requests
	unreg  unRegistered // ip (unregistered users) and its requests
	banned api.Banned   // banned users or ip
}

func NewSes() *Ses {
	return &Ses{
		reg: registered{
			uuid: make(map[string]*Session),
			id:   make(map[string]*Session),
		},
		unreg: unRegistered{
			ip: make(map[string]*requests),
		},
		banned: bann.NewBanned(),
	}
}

func (s *Ses) CheckSession(fn func(*object.CookieInfo,
	object.Status, http.ResponseWriter, *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get ip from request
		ip := getIP(r)
		// check ban_list by ip (DDoS)
		exist, expire := s.banned.UnlockBan(time.Now(), ip)
		if exist && !expire {
			fn(nil, object.ByCode(config.Code503), w, r) // block request
			return
		}
		// create object back-cookie for db
		ck := object.NewCookieInfo()
		// get userID from web-cookie
		err := ck.CookieUserIDRead(r)
		if err == http.ErrNoCookie {
			// add ip to unregistered requests
			if s.unreg.make(ip) {
				fn(object.NewCookieInfo(), nil, w, r) // start as unregistered user
			} else {
				// added to ban_list
				s.makeBan(config.UnregUser, ip)
				fn(nil, object.ByCode(config.Code503), w, r) // block request
			}
			return
		} else if err != nil { // damaged id
			fn(nil, object.ByCodeAndLog(config.Code400,
				err, "cookie: cookieUserIDRead: atoi"), w, r) // block request
			return
		}
		// check ban_list by user_id (DDoS)
		exist, expire = s.banned.UnlockBan(time.Now(), ck.UserStr)
		if exist && !expire {
			fn(nil, object.ByCode(config.Code503), w, r) // block request
			return
		}
		// get uuid-session from web-cookie
		err = ck.CookieSessionRead(r)
		if err == http.ErrNoCookie {
			// if user_id = ok but no session - clear map with this user_id
			s.reg.delete(ck.UserStr)
			fn(object.NewCookieInfo(), nil, w, r) // start as unregistered user
			return
		}
		// check user{user_id, uuid-session} in map
		if !s.reg.match(ck) {
			// delete from map
			s.reg.delete(ck.UserStr)
			s.reg.delete(ck.UUID)
			// delete cookie in browser
			object.CookieSessionAndUserID(w, []string{"", ""}, "erase")
			fn(nil, object.ByText(nil, config.AccessDenied), w, r) // block request
			return
		}
		// expired session will be deleted from map and web-cookie
		// and a new one created
		s.reg.actualizeDate(w, ck.UserStr)
		ck.Session = true
		fn(ck, nil, w, r)
	}
}

func (s *Ses) makeBan(by, key string) {
	switch by {
	case config.RegUser:
		s.reg.delete(key)
	case config.UnregUser:
		s.unreg.delete(key)
	}
	banTill := time.Now().Add(config.BanDuration)
	s.banned.Add(key, banTill)
}

type Session struct {
	id, uuid string
	expire   time.Time
	requests
}

func (s Session) expired() bool {
	return s.expire.Before(time.Now())
}

// // CheckUUID checks for registered users their session
// func (s *Ses) CheckUUID(uuid, id string, rewrite bool) bool {
// 	s.RLock()
// 	defer s.RUnlock()
// 	now := time.Now()
// 	if v, ok := s.user[uuid]; ok && v.userId == id {
// 		if v.expire.After(now) {
// 			s.Delete(uuid)
// 			if rewrite {
// 			}

// 		}
// 	}
// }

// // CheckIP checks for unregistered users for their sessions
// func (s *Ses) CheckIP(uuid, id string) {
// 	now := time.Now()
// 	s.RLock()
// 	defer s.RUnlock()
// 	if v, ok := s.user[uuid]; ok {
// 		if v.expire.After(now) {
// 		}
// 	}
// }

// func (s Session) GetUserId() string {
// 	return s.userId
// }

func genUUID() (string, object.Status) {
	key, err := uuid.NewV4()
	if err != nil {
		return "", object.ByCodeAndLog(config.Code500,
			err, "middleware create: generate uuid")
	}
	return key.String(), nil
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
