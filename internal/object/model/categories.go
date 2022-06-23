package model

import (
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Categories struct {
	Slice []*Category
	St    *object.Settings
	Ck    *object.CookieInfo
}

func NewCategories(st *object.Settings, ck *object.CookieInfo) *Categories {
	p := new(Categories)
	if st == nil {
		p.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		p.St = st
	}
	if ck == nil {
		p.Ck = new(object.CookieInfo)
	} else {
		p.Ck = ck
	}
	return p
}

func (c *Categories) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.St.Key[key] = data
	} else {
		c.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (c *Categories) GetList() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := c.St.Key[config.KeyPost]; ok {
		qs.QueryName = config.QueSelectCategoryBy
		qs.QueryFields = []interface{}{
			config.FieldPost,
		}
		if value == nil {
			qs.Fields = []interface{}{0}
		} else {
			qs.Fields = value
		}
		return qs
	}
	if value, ok := c.St.Key[config.KeyID]; ok {
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
		return qs
	}
	return &object.QuerySettings{
		QueryName: config.QueSelectCategories,
	}
}

func (c *Categories) NewList() []interface{} {
	category := new(Category)
	c.Slice = append(c.Slice, category)
	if _, ok := c.St.Key[config.KeyID]; ok {
		return []interface{}{
			&category.ID,
		}
	}
	return []interface{}{
		&category.ID,
		&category.Name,
	}
}

func (c *Categories) IfNil() interface{} {
	return []*Category{new(Category).ifNil()}
}

func (c *Categories) Return() *Buf {
	return &Buf{Categories: c}
}
