package model

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Post struct {
	ID         int
	Title      string
	Body       string
	User       int
	Name       string
	Created    time.Time
	Categories any
	Likes      any
	Comments   *Comments
	Liked      any
	Image      any
	St         *object.Settings
	Ck         *object.CookieInfo
}

func NewPost(st *object.Settings, ck *object.CookieInfo) *Post {
	p := new(Post)
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

func (p *Post) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		p.St.Key[key] = data
	} else {
		p.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (p *Post) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueSelectPostsBy
	qs.QueryFields = []interface{}{
		config.TabPosts,
		config.FieldID,
	}
	if value, ok := p.St.Key[config.KeyPost]; ok {
		qs.Fields = value
	} else {
		qs.Fields = []interface{}{0} // for null list
	}
	return qs
}

func (p *Post) New() []interface{} {
	return []interface{}{
		&p.ID,
		&p.Title,
		&p.Body,
		&p.User,
		&p.Name,
		&p.Created,
		&p.Image,
	}
}

func (p *Post) IfNil() interface{} {
	return p.ifNil()
}

func (p *Post) ifNil() *Post {
	return &Post{
		Title:   "no posts created",
		Body:    "sorry, empty here",
		Created: time.Now(),
		User:    1,
		Name:    "Admin",
	}
}
