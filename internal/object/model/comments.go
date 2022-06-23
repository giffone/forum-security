package model

import (
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Comments struct {
	Slice []*Comment
	St    *object.Settings
	Ck    *object.CookieInfo
}

func NewComments(st *object.Settings, ck *object.CookieInfo) *Comments {
	c := new(Comments)
	if st == nil {
		c.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		c.St = st
	}
	if ck == nil {
		c.Ck = new(object.CookieInfo)
	} else {
		c.Ck = ck
	}
	return c
}

func (c *Comments) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		c.St.Key[key] = data
	} else {
		c.St.Key[config.KeyComment] = []interface{}{0}
	}
}

func (c *Comments) GetList() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := c.St.Key[config.KeyPost]; ok {
		qs.QueryName = config.QueSelectCommentsBy
		qs.QueryFields = []interface{}{
			config.TabComments,
			config.FieldPost,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.Post}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[config.KeyUser]; ok {
		qs.QueryName = config.QueSelectCommentsBy
		qs.QueryFields = []interface{}{
			config.TabComments,
			config.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[config.KeyRated]; ok {
		qs.QueryName = config.QueSelectCommentsRatedBy
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{c.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := c.St.Key[config.KeyID]; ok {
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
	}
	return qs
}

func (c *Comments) NewList() []interface{} {
	comment := new(Comment)
	c.Slice = append(c.Slice, comment)
	if _, ok := c.St.Key[config.KeyRated]; ok {
		return []interface{}{
			&comment.ID,
			&comment.Name,
			&comment.Body,
			&comment.Created,
			&comment.Post,
			&comment.Liked,
		}
	} else if _, ok := c.St.Key[config.KeyID]; ok {
		return []interface{}{
			&comment.ID,
		}
	}
	return []interface{}{
		&comment.ID,
		&comment.Name,
		&comment.Body,
		&comment.Created,
		&comment.Post,
	}
}

func (c *Comments) IfNil() interface{} {
	return []*Comment{new(Comment).ifNil()}
}

func (c *Comments) Return() *Buf {
	return &Buf{Comments: c}
}

func (c *Comments) LSlice() int {
	return len(c.Slice)
}

func (c *Comments) PostOrCommentID(index int) int {
	return c.Slice[index].ID
}

func (c *Comments) Add(key string, index int, data interface{}) {
	switch key {
	case config.KeyLike:
		c.Slice[index].Likes = data
	case config.KeyRated:
		c.Slice[index].Liked = data
	}
}

func (c *Comments) Cookie() *object.CookieInfo {
	return c.Ck
}

func (c *Comments) Settings() *object.Settings {
	return c.St
}

func (c *Comments) KeyRole() string {
	return config.KeyComment
}

func (c *Comments) KeyLiked() string {
	return config.KeyCommentRated
}
