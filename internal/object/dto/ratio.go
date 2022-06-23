package dto

import (
	"net/http"
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

const (
	like    = 1
	dislike = 2
)

type Ratio struct {
	PostOrComm map[string]string // post or comment - map[post]id
	Ratio      int               // like/dislike
	Obj        object.Obj
}

func NewRatio(st *object.Settings, sts *object.Statuses, ck *object.CookieInfo) *Ratio {
	l := new(Ratio)
	l.Obj.NewObjects(st, sts, ck)
	return l
}

func (r *Ratio) AddByPOST(request *http.Request) bool {
	// ratio - like/dislike
	rate := request.PostFormValue(config.KeyRate)
	if rate == config.KeyLike {
		r.Ratio = like
	} else if rate == config.KeyDislike {
		r.Ratio = dislike
	} else {
		r.Obj.Sts.ByCodeAndLog(config.Code400,
			nil, "dto: like: postFormValue is nil")
		return false
	}
	// post or comment ID
	postOrComm := config.KeyPost
	objID := request.PostFormValue(postOrComm)
	if objID == "" {
		postOrComm = config.KeyComment
		objID = request.PostFormValue(postOrComm)
	}
	if objID == "" {
		r.Obj.Sts.ByCodeAndLog(config.Code400,
			nil, "dto: like: post or comment id is nil")
		return false
	}
	r.PostOrComm = make(map[string]string)
	r.PostOrComm[postOrComm] = objID
	return true
}

func (r *Ratio) AddByGET(request *http.Request) bool {
	// get post id
	sts := r.Obj.Ck.CookiePostIDRead(request)
	if sts != nil {
		r.Obj.Sts = sts.Status()
		return false
	}
	// read url
	u := request.URL.Query()
	// ratio - like/dislike
	rate := u.Get(config.KeyRate)
	if rate == config.KeyLike {
		r.Ratio = like
	} else if rate == config.KeyDislike {
		r.Ratio = dislike
	} else {
		r.Obj.Sts.ByCodeAndLog(config.Code400,
			nil, "dto: like: postFormValue is nil")
		return false
	}
	// post or comment ID
	postOrComm := config.KeyPost
	objID := u.Get(postOrComm)
	if objID == "" {
		postOrComm = config.KeyComment
		objID = u.Get(postOrComm)
	}
	if objID == "" {
		r.Obj.Sts.ByCodeAndLog(config.Code400,
			nil, "dto: like: post or comment id is nil")
		return false
	}
	r.PostOrComm = make(map[string]string)
	r.PostOrComm[postOrComm] = objID
	return true
}

func (r *Ratio) MakeKeys(key string, data ...interface{}) {
	if key != "" {
		r.Obj.St.Key[key] = data
	} else {
		r.Obj.St.Key[config.KeyPost] = []interface{}{0}
	}
}

// Create prepares query for db and choose fields for adding incoming data
func (r *Ratio) Create() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueInsert4
	if _, ok := r.Obj.St.Key[config.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.FieldUser,
			config.FieldPost,
			config.FieldLike,
			config.FieldCreated,
		}
		qs.Fields = []interface{}{
			r.Obj.Ck.User,
			r.PostOrComm[config.KeyPost],
			r.Ratio,
			time.Now(),
		}
	} else {
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.FieldUser,
			config.FieldComment,
			config.FieldLike,
			config.FieldCreated,
		}
		qs.Fields = []interface{}{
			r.Obj.Ck.User,
			r.PostOrComm[config.KeyComment],
			r.Ratio,
			time.Now(),
		}
	}
	return qs
}

func (r *Ratio) Delete() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueDeleteBy
	if value, ok := r.Obj.St.Key[config.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			config.TabPostsLikes,
			config.TabPostsLikes,
			config.FieldID,
		}
		qs.Fields = value
	} else if value, ok := r.Obj.St.Key[config.KeyComment]; ok {
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.FieldID,
		}
		qs.Fields = value // for null list
	} else {
		qs.QueryFields = []interface{}{
			config.TabCommentsLikes,
			config.TabCommentsLikes,
			config.FieldID,
		}
		qs.Fields = []interface{}{0} // for null list
	}
	return qs
}
