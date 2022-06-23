package account

import (
	"context"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/object/model"
	"github.com/giffone/forum-security/internal/service"
)

type hAccount struct {
	sPost     service.Post
	sCategory service.Category
	sComment  service.Comment
	sRatio    service.Ratio
}

func NewHandler(sPost service.Post, sCategory service.Category,
	sComment service.Comment, sRatio service.Ratio,
) api.Handler {
	return &hAccount{
		sPost:     sPost,
		sCategory: sCategory,
		sComment:  sComment,
		sRatio:    sRatio,
	}
}

func (ha *hAccount) Register(ctx context.Context, router *http.ServeMux, session api.Middleware) {
	router.HandleFunc(config.URLAccount, session.CheckSession(ctx, ha.ByUser))
	router.HandleFunc(config.URLAccountRatio, session.CheckSession(ctx, ha.CreateLike))
}

func (ha *hAccount) ByUser(ctx context.Context, ck *object.CookieInfo, sts object.Status,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	// clear unused postID
	object.CookiePostIDDel(w)
	// get data from db, parse and execute response
	ha.get(ctx, ck, w)
}

func (ha *hAccount) CreateLike(ctx context.Context, ck *object.CookieInfo, sts object.Status,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	ctx, cancel := context.WithTimeout(ctx, config.TimeLimit10s)
	defer cancel()
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "POST" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	// need session always to continue
	if !ck.Session {
		api.Message(w, object.ByText(nil, config.AccessDenied))
		return
	}
	// create DTO with a new rate
	ratio := dto.NewRatio(nil, nil, ck)
	// add request data to DTO and check err
	if !ratio.AddByPOST(r) {
		api.Message(w, ratio.Obj.Sts)
		return
	}
	// make ratio
	_, sts = ha.sRatio.Create(ctx, ratio)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	redirect := r.PostFormValue(config.KeyLink)
	http.Redirect(w, r, config.URLAccount+redirect, config.Code302)
	// ha.get(ctx, ck, w)
}

func (ha *hAccount) get(ctx context.Context, ck *object.CookieInfo, w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(ctx, config.TimeLimit10s)
	defer cancel()
	// parse
	pe, sts := api.NewParseExecute("account").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	//
	pe.Data["Acc"] = true
	// link for "form action" submit
	pe.Data["RatioLink"] = config.URLAccountRatio
	// insert session
	pe.Data["Session"] = ck.Session
	// create new model posts
	posts := model.NewPosts(nil, ck)
	// make keys for sort posts
	posts.MakeKeys(config.KeyUser)
	// insert posts
	pe.Data["Posts"], sts = ha.sPost.Get(ctx, posts)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model posts
	postsRated := model.NewPosts(nil, ck)
	// make keys for sort posts
	postsRated.MakeKeys(config.KeyRated)
	// insert posts
	pe.Data["PostsRated"], sts = ha.sPost.Get(ctx, postsRated)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// need make refer to post
	st := object.NewSettings()
	st.Refers = true
	// create new model comments
	comments := model.NewComments(st, ck)
	// make keys for sort posts
	comments.MakeKeys(config.KeyUser)
	// insert posts
	pe.Data["Comments"], sts = ha.sComment.Get(ctx, comments)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model comments
	commentsRated := model.NewComments(st.ClearKey(), ck)
	// make keys for sort posts
	commentsRated.MakeKeys(config.KeyRated)
	// insert posts
	pe.Data["CommentsRated"], sts = ha.sComment.Get(ctx, commentsRated)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model categories
	c := model.NewCategories(nil, nil)
	// insert categories
	pe.Data["Category"], sts = ha.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// insert method to show - one post or all posts
	pe.Data["AllPost"] = posts.St.AllPost
	// execute
	pe.Execute(w, config.Code200)
}
