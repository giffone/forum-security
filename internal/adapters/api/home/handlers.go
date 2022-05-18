package home

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/model"
	"github.com/giffone/forum-security/internal/service"
)

type hHome struct {
	sPost     service.Post
	sCategory service.Category
}

func NewHandler(sPost service.Post, sCategory service.Category) api.Handler {
	return &hHome{
		sPost:     sPost,
		sCategory: sCategory,
	}
}

func (hh *hHome) Register(ctx context.Context, router *http.ServeMux, session api.Middleware) {
	router.HandleFunc(constant.URLHome, session.CheckSession(ctx, hh.Home))
	router.HandleFunc(constant.URLFavicon, hh.Favicon)
	router.HandleFunc(constant.URLCategoryBy, session.CheckSession(ctx, hh.ByCategory))
}

func (hh *hHome) Home(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.ByCode(constant.Code405))
		return
	}
	if r.URL.Path != "/" {
		api.Message(w, object.ByCode(constant.Code404))
		return
	}
	posts := model.NewPosts(nil, ck)
	posts.St.AllPost = true
	hh.get(ctx, posts, w)
}

func (hh *hHome) ByCategory(ctx context.Context, ck *object.Cookie, sts object.Status,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	// check errors in cookie
	if sts != nil {
		api.Message(w, sts)
		return
	}
	if r.Method != "GET" {
		api.Message(w, object.ByCode(constant.Code405))
		return
	}
	// get id category from url
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, constant.URLCategoryBy))
	if err != nil || id == 0 {
		sts := object.ByCodeAndLog(constant.Code400,
			err, "handler: postFormValue: atoi")
		api.Message(w, sts)
		return
	}
	posts := model.NewPosts(nil, ck)
	posts.MakeKeys(constant.KeyCategory, id)
	posts.St.AllPost = true
	hh.get(ctx, posts, w)
}

func (hh *hHome) Favicon(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " ", r.URL.Path)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "assets/ico/favicon.ico")
}

func (hh *hHome) get(ctx context.Context, m model.Models, w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(ctx, constant.TimeLimit)
	defer cancel()
	// parse
	pe, sts := api.NewParseExecute("index").Parse()
	if sts != nil {
		api.Message(w, sts)
		return
	}
	p := m.Return().Posts
	// session
	pe.Data["Session"] = p.Ck.Session
	// get data
	pe.Data["AllPost"] = p.St.AllPost
	pe.Data["Posts"], sts = hh.sPost.Get(ctx, m)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	c := model.NewCategories(nil, nil)
	pe.Data["Category"], sts = hh.sCategory.GetList(ctx, c)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// execute
	pe.Execute(w, constant.Code200)
}
