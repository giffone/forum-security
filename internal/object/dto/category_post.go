package dto

import (
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
)

type CategoryPost struct {
	Post     int      // current new created
	Category int      // current id from slice
	Slice    []string // from http response
	ID       []int    // checked for valid from http response
}

func NewCategoryPost() *CategoryPost {
	return new(CategoryPost)
}

func (cp *CategoryPost) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert2,
		QueryFields: []interface{}{constant.TabPostsCategories,
			constant.FieldPost, constant.FieldCategory},
		Fields: []interface{}{
			cp.Post,
			cp.Category,
		},
	}
}

func (cp *CategoryPost) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}
