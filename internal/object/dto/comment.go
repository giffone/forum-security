package dto

import (
	"net/http"
	"strings"
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Comment struct {
	Body string
	Obj  object.Obj
}

func NewComment(st *object.Settings, sts *object.Statuses, ck *object.CookieInfo) *Comment {
	c := new(Comment)
	c.Obj.NewObjects(st, sts, ck)
	return c
}

func (c *Comment) Add(r *http.Request) bool {
	// get user id
	sts := c.Obj.Ck.CookieUserIDRead(r)
	if sts != nil {
		// c.Obj.Sts = sts.Status()                   /////////////////// here error
		return false
	}
	// get post id
	// sts = c.Obj.Ck.CookiePostIDRead(r)                   /////////////////// here error
	if sts != nil {
		// c.Obj.Sts = sts.Status()                   /////////////////// here error
		return false
	}
	c.Body = r.PostFormValue("body text")
	return true
}

func (c *Comment) Valid() bool {
	// delete space for check an any symbol
	body := strings.TrimSpace(c.Body)
	if body == "" {
		c.Obj.Sts.ByText(nil, config.TooShort,
			"comment", "one")
		return false
	}
	return true
}

func (c *Comment) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: config.QueInsert4,
		QueryFields: []interface{}{
			config.TabComments,
			config.FieldUser,
			config.FieldPost,
			config.FieldBody,
			config.FieldCreated,
		},
		Fields: []interface{}{
			c.Obj.Ck.User,
			c.Obj.Ck.Post,
			c.Body,
			time.Now(),
		},
	}
}

func (c *Comment) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}
