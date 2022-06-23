package model

import "github.com/giffone/forum-security/internal/object"

type Model interface {
	New() []interface{}
	Get() *object.QuerySettings
}

type Models interface {
	NewList() []interface{}
	GetList() *object.QuerySettings
	Return() *Buf
}

// Buf for return struct by name
type Buf struct {
	Category   *Category
	Categories *Categories
	Comment    *Comment
	Comments   *Comments
	Like       *Ratio
	LCount     *LCount
	LikesCount *LikesCount
	Post       *Post
	Posts      *Posts
	Session    *Session
	User       *User
	Users      *Users
	CheckID    *CheckID
}

type PostOrComment interface {
	LSlice() int
	PostOrCommentID(index int) int
	Add(key string, index int, data interface{})
	Cookie() *object.CookieInfo
	Settings() *object.Settings
	KeyRole() string
	KeyLiked() string
}
