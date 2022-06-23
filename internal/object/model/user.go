package model

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type User struct {
	ID       int
	Login    string
	Name     string
	Password string
	Email    string
	Root     int
	Created  time.Time
	St       *object.Settings
	Ck       *object.CookieInfo
}

func NewUser(st *object.Settings, ck *object.CookieInfo) *User {
	u := new(User)
	if st == nil {
		u.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		u.St = st
	}
	if ck == nil {
		u.Ck = new(object.CookieInfo)
	} else {
		u.Ck = ck
	}
	return u
}

func (u *User) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		u.St.Key[key] = data
	} else {
		u.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (u *User) Get() *object.QuerySettings {
	return get(u.St.Key)
}

func (u *User) New() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Login,
		&u.Name,
		&u.Password,
		&u.Email,
		&u.Root,
		&u.Created,
	}
}

func get(key map[string][]interface{}) *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueSelectUserBy
	if value, ok := key[config.KeyID]; ok {
		qs.QueryFields = []interface{}{
			config.FieldID,
		}
		qs.Fields = value
	} else if value, ok := key[config.KeyLogin]; ok {
		qs.QueryFields = []interface{}{
			config.FieldLogin,
		}
		qs.Fields = value
	}
	return qs
}
