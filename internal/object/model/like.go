package model

import (
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type Ratio struct {
	PostOrComm int       // current post_id or comment_id
	Ratio      int       // like_id or dislike_id
	Body       string    // like or dislike (name)
	Created    time.Time // created date
	St         *object.Settings
	Ck         *object.CookieInfo
}

func NewLike(st *object.Settings, ck *object.CookieInfo) *Ratio {
	l := new(Ratio)
	if st == nil {
		l.St = &object.Settings{
			Key: make(map[string][]interface{}),
		}
	} else {
		l.St = st
	}
	if ck == nil {
		l.Ck = new(object.CookieInfo)
	} else {
		l.Ck = ck
	}
	return l
}

func (l *Ratio) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		l.St.Key[key] = data
	} else {
		l.St.Key[config.KeyPost] = []interface{}{0}
	}
}

func (l *Ratio) Get() *object.QuerySettings {
	qs := new(object.QuerySettings)
	if value, ok := l.St.Key[config.KeyPost]; ok {
		qs.QueryName = config.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.FieldUser,
			config.TabPostsLikes,
			config.FieldPost,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[config.KeyComment]; ok {
		qs.QueryName = config.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.FieldUser,
			config.TabCommentsLikes,
			config.FieldComment,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[config.KeyPostRated]; ok {
		qs.QueryName = config.QueSelectLikedOrNot
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.FieldUser,
			config.TabPostsLikes,
			config.FieldPost,
		}
		qs.Fields = value
	} else if value, ok := l.St.Key[config.KeyCommentRated]; ok {
		qs.QueryName = config.QueSelectLikedOrNot
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.FieldUser,
			config.TabCommentsLikes,
			config.FieldComment,
		}
		qs.Fields = value
	} else { // for null list
		qs.QueryName = config.QueSelectLikeBy
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.FieldUser,
			config.TabPostsLikes,
			config.FieldPost,
		}
		qs.Fields = []interface{}{0, 0}
	}
	return qs
}

func (l *Ratio) New() []interface{} {
	if _, ok := l.St.Key[config.KeyPostRated]; ok {
		return []interface{}{
			&l.Body,
		}
	} else if _, ok := l.St.Key[config.KeyCommentRated]; ok {
		return []interface{}{
			&l.Body,
		}
	}
	return []interface{}{
		&l.PostOrComm,
		&l.Ratio,
		&l.Body,
		&l.Created,
	}
}
