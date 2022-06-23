package object

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/giffone/forum-security/internal/config"
)

type CookieInfo struct {
	User       int
	UserStr    string
	Post       int
	PostString string
	Session    bool
	UUID       string
}

func NewCookieInfo() *CookieInfo {
	return new(CookieInfo)
}

func (ci *CookieInfo) CookieUserIDRead(r *http.Request) error {
	ck, err := r.Cookie(config.CookieUserID)
	if err != nil {
		return err
	}
	user, err := strconv.Atoi(ck.Value)
	if err != nil || user == 0 {
		return err
	}
	ci.User = user
	ci.UserStr = ck.Value
	log.Printf("object cookie: userID is: %d", ci.User)
	return nil
}

func (ci *CookieInfo) CookieSessionRead(r *http.Request) error {
	ck, err := r.Cookie(config.CookieSession)
	if err != nil {
		return err
	}
	ci.UUID = ck.Value
	log.Printf("object cookie: middleware-uuid is: %s", ci.UUID)
	return nil
}

func (ci *CookieInfo) CookiePostIDRead(r *http.Request) Status {
	p, err := r.Cookie(config.CookiePostID)
	if err != nil {
		return ByCodeAndLog(config.Code400,
			err, "cookie: cookiePostIDRead")
	}
	ci.PostString = p.Value
	log.Printf("object cookie: postID is: %s", ci.PostString)
	return nil
}

func (ci *CookieInfo) AddUser(id int) *CookieInfo {
	ci.User = id
	return ci
}

func CookieSessionAndUserID(w http.ResponseWriter, value []string, method string) Status { /////////////////// here error
	// name := []string{config.CookieSession, config.CookieUserID}
	// sts := cookieSet(w, name, value, method)
	// if sts != nil {
	// 	return sts
	// }
	return nil
}

func CookiePostID(w http.ResponseWriter, id string) Status { /////////////////// here error
	// name := []string{config.CookiePostID}
	// value := []string{id}
	// sts := cookieSet(w, name, value, "")
	// if sts != nil {
	// 	return sts
	// }
	return nil
}

func CookiePostIDDel(w http.ResponseWriter) Status { /////////////////// here error
	// name := []string{config.CookiePostID}
	// value := []string{""}
	// sts := cookieSet(w, name, value, "erase")
	// if sts != nil {
	// 	return sts
	// }
	return nil
}

type Cookie struct {
	src  map[string]string
	save bool
	w    http.ResponseWriter
}

func NewCookie(w http.ResponseWriter, src map[string]string) *Cookie {
	return &Cookie{src: src, w: w}
}

// Set sets cookie into ResponseWriter
// true - save, false - delete cookie
func (c *Cookie) Set(saveOrDelete bool) {
	c.save = saveOrDelete
	c.set()
}

func (c Cookie) set() {
	for key, value := range c.src {
		ck := &http.Cookie{}
		ck.Name = key
		ck.Value = value
		ck.Path = "/"
		if c.save {
			ck.Expires = time.Now().Add(config.SessionExpire)
			ck.MaxAge = config.SessionMaxAge
		} else {
			ck.Expires = time.Unix(1, 0)
			ck.MaxAge = -1
		}
		http.SetCookie(c.w, ck)
	}
}
