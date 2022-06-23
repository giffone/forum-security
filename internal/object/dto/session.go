package dto

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Session struct {
	Expire time.Time
	Obj    object.Obj
}

func NewSession(st *object.Settings, sts *object.Statuses, ck *object.CookieInfo) *Session {
	s := new(Session)
	s.Obj.NewObjects(st, sts, ck)
	return s
}

func (s *Session) Add(date time.Time) {
	s.Expire = date
}

func (s *Session) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: config.QueInsert3,
		QueryFields: []interface{}{
			config.TabSessions,
			config.FieldUser,
			config.FieldUUID,
			config.FieldExpire,
		},
		Fields: []interface{}{
			s.Obj.Ck.User,
			s.Obj.Ck.UUID,
			s.Expire,
		},
	}
}

func (s *Session) Delete() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: config.QueDeleteBy,
		QueryFields: []interface{}{
			config.TabSessions,
			config.TabSessions,
			config.FieldUser,
		},
		Fields: []interface{}{
			s.Obj.Ck.User,
		},
	}
}
