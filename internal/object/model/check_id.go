package model

import (
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type CheckID struct {
	ID  int
	Obj object.Obj
}

func NewCheckID(st *object.Settings, sts *object.Statuses, ck *object.CookieInfo) *CheckID {
	c := new(CheckID)
	c.Obj.NewObjects(st, sts, ck)
	return c
}

func (c *CheckID) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.Obj.St.Key[key] = data
	} else {
		c.Obj.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (c *CheckID) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueSelect
	if value, ok := c.Obj.St.Key[config.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			config.TabPosts,
			config.TabPosts,
			config.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[config.KeyCategory]; ok {
		qs.QueryName = config.QueSelect
		qs.QueryFields = []interface{}{
			config.TabCategories,
			config.TabCategories,
			config.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[config.KeyComment]; ok {
		qs.QueryName = config.QueSelect
		qs.QueryFields = []interface{}{
			config.TabComments,
			config.TabComments,
			config.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[config.KeyLogin]; ok {
		qs.QueryName = config.QueSelect
		qs.QueryFields = []interface{}{
			config.TabUsers,
			config.TabUsers,
			config.FieldLogin,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.Obj.St.Key[config.KeyEmail]; ok {
		qs.QueryName = config.QueSelect
		qs.QueryFields = []interface{}{
			config.TabUsers,
			config.TabUsers,
			config.FieldEmail,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
	}
	return qs
}

func (c *CheckID) New() []interface{} {
	return []interface{}{
		&c.ID,
	}
}
