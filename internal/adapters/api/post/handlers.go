package post

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/object/model"
	"github.com/giffone/forum-security/internal/service"
)

type hPost struct {
	service     service.Post
	sCategory   service.Category
	sComment    service.Comment
	sRatio      service.Ratio
	sMiddleware service.Middleware
}

func NewHandler(service service.Post,
	sCategory service.Category, sComment service.Comment,
	sRatio service.Ratio, sMiddleware service.Middleware,
) api.Handler {
	return &hPost{
		service:     service,
		sCategory:   sCategory,
		sComment:    sComment,
		sRatio:      sRatio,
		sMiddleware: sMiddleware,
	}
}

func (hp *hPost) Register(ctx context.Context, router *http.ServeMux, middleware api.Middleware) {
	router.HandleFunc(config.URLRead, middleware.CheckSession(ctx, hp.Read))
	router.HandleFunc(config.URLPost, middleware.CheckSession(ctx, hp.CreatePost))
	router.HandleFunc(config.URLComment, middleware.CheckSession(ctx, hp.CreateComment))
	router.HandleFunc(config.URLReadRatio, middleware.CheckSession(ctx, hp.CreateRatio))
}

func (hp *hPost) Read(ctx context.Context, ck *object.CookieInfo, sts object.Status,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	// git post id
	ck.PostString = r.URL.Query().Get("post")
	// check valid id for refer page
	post := dto.NewCheckIDAtoi(config.KeyPost, ck.PostString)
	idWho, sts := hp.sMiddleware.GetID(ctx, post)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// save id
	ck.Post = idWho
	// save id in cookie
	object.CookiePostID(w, ck.PostString)
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreatePost(ctx context.Context, ck *object.CookieInfo, sts object.Status,
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
	// limit to read 32Mb
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
	// create DTO with a new post
	post := dto.NewPost(nil, nil, ck)
	// add request data to DTO & check fields for valid
	if !post.Add(r) || !post.Valid() {
		api.Message(w, post.Obj.Sts)
		return
	}
	// create post in database
	id, sts := hp.service.Create(ctx, post)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	ck.Post = id
	// send new id to cookie
	object.CookiePostID(w, strconv.Itoa(id))
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreateComment(ctx context.Context, ck *object.CookieInfo, sts object.Status,
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
	// create DTO with a new comment
	comment := dto.NewComment(nil, nil, ck)
	// add request data to DTO & check fields for valid
	if !comment.Add(r) || !comment.Valid() {
		api.Message(w, comment.Obj.Sts)
		return
	}
	// create comment in database
	_, sts = hp.sComment.Create(ctx, comment)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) CreateRatio(ctx context.Context, ck *object.CookieInfo, sts object.Status,
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
	if r.Method != "GET" {
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
	if !ratio.AddByGET(r) {
		api.Message(w, ratio.Obj.Sts)
		return
	}
	// create like
	_, sts = hp.sRatio.Create(ctx, ratio)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// check valid id for refer page
	postID := dto.NewCheckIDAtoi(config.KeyPost, ck.PostString)
	idWho, sts := hp.sMiddleware.GetID(ctx, postID)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// save id
	ck.Post = idWho
	// get data from db, parse and execute response
	hp.get(ctx, ck, w)
}

func (hp *hPost) get(ctx context.Context, ck *object.CookieInfo, w http.ResponseWriter) {
	// parse
	pe, sts := api.NewParseExecute("post").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// link for "form action" submit
	pe.Data["RatioLink"] = config.URLReadRatio
	// insert session
	pe.Data["Session"] = ck.Session
	// create new model posts
	p := model.NewPosts(nil, ck)
	// make keys for sort posts
	p.MakeKeys(config.KeyPost)
	// insert posts
	pe.Data["Posts"], sts = hp.service.Get(ctx, p)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// insert method to show - one post or all posts
	pe.Data["AllPost"] = p.St.AllPost
	// create new model categories
	c := model.NewCategories(nil, nil)
	// insert categories
	pe.Data["Category"], sts = hp.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// create new model comments
	cm := model.NewComments(nil, ck)
	// make keys for sort posts
	cm.MakeKeys(config.KeyPost)
	// insert comments
	pe.Data["Comments"], sts = hp.sComment.Get(ctx, cm)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// execute
	pe.Execute(w, config.Code200)
}
