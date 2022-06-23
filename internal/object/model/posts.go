package model

import (
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Posts struct {
	Slice []*Post
	St    *object.Settings
	Ck    *object.CookieInfo
}

func NewPosts(st *object.Settings, ck *object.CookieInfo) *Posts {
	p := new(Posts)
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

func (p *Posts) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		p.St.Key[key] = data
	} else {
		p.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (p *Posts) GetList() *object.QuerySettings {
	if len(p.St.Key) == 0 {
		return &object.QuerySettings{
			QueryName: config.QueSelectPosts,
		}
	}
	qs := new(object.QuerySettings)
	if value, ok := p.St.Key[config.KeyCategory]; ok {
		qs.QueryName = config.QueSelectPostsAndCategoryBy
		qs.QueryFields = []interface{}{
			config.TabPostsCategories,
			config.FieldCategory,
		}
		qs.Fields = value
	} else if value, ok := p.St.Key[config.KeyPost]; ok {
		qs.QueryName = config.QueSelectPostsBy
		qs.QueryFields = []interface{}{
			config.TabPosts,
			config.FieldID,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.Post}
		} else {
			qs.Fields = value
		}
	} else if value, ok := p.St.Key[config.KeyUser]; ok {
		qs.QueryName = config.QueSelectPostsBy
		qs.QueryFields = []interface{}{
			config.TabPosts,
			config.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.User}
		} else {
			qs.Fields = value
		}
	} else if value, ok := p.St.Key[config.KeyID]; ok {
		qs.QueryName = config.QueSelect
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
	} else if value, ok := p.St.Key[config.KeyRated]; ok {
		qs.QueryName = config.QueSelectPostsRatedBy
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.FieldUser,
		}
		if value == nil {
			qs.Fields = []interface{}{p.Ck.User}
		} else {
			qs.Fields = value
		}
	}
	return qs
}

func (p *Posts) NewList() []interface{} {
	post := new(Post)
	p.Slice = append(p.Slice, post)
	// for account handler
	if _, ok := p.St.Key[config.KeyComment]; ok {
		comment := new(Comment)
		comments := new(Comments)
		comments.Slice = append(comments.Slice, comment)
		post.Comments = comments
		return []interface{}{
			&post.ID,
			&post.Title,
			&post.Body,
			&post.User,
			&post.Name,
			&post.Created,
			&post.Image,
			&comment.ID,
			&comment.Name,
			&comment.Body,
			&comment.Created,
		}
	} else if _, ok := p.St.Key[config.KeyRated]; ok {
		return []interface{}{
			&post.ID,
			&post.Title,
			&post.Body,
			&post.User,
			&post.Name,
			&post.Created,
			&post.Image,
			&post.Liked,
		}
	} else if _, ok := p.St.Key[config.KeyID]; ok {
		return []interface{}{
			&post.ID,
		}
	}
	return []interface{}{
		&post.ID,
		&post.Title,
		&post.Body,
		&post.User,
		&post.Name,
		&post.Created,
		&post.Image,
	}
}

func (p *Posts) IfNil() interface{} {
	return []*Post{new(Post).ifNil()}
}

func (p *Posts) Return() *Buf {
	return &Buf{Posts: p}
}

func (p *Posts) LSlice() int {
	return len(p.Slice)
}

func (p *Posts) PostOrCommentID(index int) int {
	return p.Slice[index].ID
}

// Add adding information to slice post/comment by index
func (p *Posts) Add(key string, index int, data interface{}) {
	switch key {
	case config.KeyLike:
		p.Slice[index].Likes = data
	case config.KeyRated:
		p.Slice[index].Liked = data
	case config.KeyCategory:
		p.Slice[index].Categories = data
	}
}

func (p *Posts) Cookie() *object.CookieInfo {
	return p.Ck
}

func (p *Posts) Settings() *object.Settings {
	return p.St
}

func (p *Posts) KeyRole() string {
	return config.KeyPost
}

func (p *Posts) KeyLiked() string {
	return config.KeyPostRated
}
