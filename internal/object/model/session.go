package model

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Session struct {
	User   string
	UUID   string
	Expire time.Time
	St     *object.Settings // settings
	Ck     *object.CookieInfo
}

func NewSession(st *object.Settings, ck *object.CookieInfo) *Session {
	s := new(Session)
	if st == nil {
		s.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		s.St = st
	}
	if ck == nil {
		s.Ck = new(object.CookieInfo)
	} else {
		s.Ck = ck
	}
	return s
}

func (s *Session) Get() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: config.QueSelectSessionBy,
		QueryFields: []interface{}{
			config.FieldUser,
			config.FieldUUID,
		},
		Fields: []interface{}{
			s.Ck.User,
			s.Ck.UUID,
		},
	}
}

func (s *Session) New() []interface{} {
	return []interface{}{
		&s.User,
		&s.UUID,
		&s.Expire,
	}
}
