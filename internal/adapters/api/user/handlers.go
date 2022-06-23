package user

import (
	"context"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/service"
)

type hUser struct {
	service service.User
	auth    *config.Auth
}

func NewHandler(service service.User, auth *config.Auth) api.Handler {
	return &hUser{
		service: service,
		auth:    auth,
	}
}

func (hu *hUser) Register(ctx context.Context, router *http.ServeMux, s api.Middleware) {
	router.HandleFunc(config.URLSignUp, s.Skip(ctx, hu.SignUp))
	router.HandleFunc(config.URLLogin, s.Skip(ctx, hu.Login))
	router.HandleFunc(config.URLLogout, s.Skip(ctx, hu.Logout))
}

func (hu *hUser) SignUp(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method != "POST" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	ctx, cancel := context.WithTimeout(ctx, config.TimeLimit10s)
	defer cancel()

	// create DTO with a new user
	user := dto.NewUser(nil, nil)
	// create return page
	user.Obj.Sts.ReturnPage = config.URLLogin + "?#signup"
	// add data from request
	user.Add(r)
	// and check fields for incorrect data entry
	if !user.ValidLogin() || !user.ValidPassword() ||
		!user.ValidEmail() || !user.CryptPassword() {
		api.Message(w, user.Obj.Sts)
		return
	}
	// create user in database
	id, sts := hu.service.Create(ctx, user)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// make session
	method := ""
	if m := r.PostFormValue("remember"); m == "on" {
		method = "remember"
	}
	sts = ses.CreateSession(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.ByText(nil, config.StatusCreated,
		"to return on main page click button below")
	api.Message(w, sts)
}

func (hu *hUser) Login(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request,
) {
	log.Println(r.Method, " ", r.URL.Path)
	if r.Method == "GET" {
		pe, sts := api.NewParseExecute("login").Parse()
		if sts != nil {
			api.Message(w, sts)
			return
		}
		// link for refers login
		if !hu.auth.Github.Empty {
			pe.Data["Github"] = config.URLLoginGithub
		}
		if !hu.auth.Facebook.Empty {
			pe.Data["Facebook"] = config.URLLoginFacebook
		}
		if !hu.auth.Google.Empty {
			pe.Data["Google"] = config.URLLoginGoogle
		}
		pe.Execute(w, config.Code200)
		return
	}
	if r.Method != "POST" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	ctx, cancel := context.WithTimeout(ctx, config.TimeLimit10s)
	defer cancel()

	// create DTO with a user
	user := dto.NewUser(nil, nil)
	// create return page
	// must be before Add() for ignore re-password check
	user.Obj.Sts.ReturnPage = config.URLLogin
	// add data from request
	user.Add(r)
	// and check fields for incorrect data entry
	if !user.ValidLogin() || !user.ValidPassword() {
		api.Message(w, user.Obj.Sts)
		return
	}
	// checks login password
	id, sts := hu.service.CheckLoginPassword(ctx, user)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// make session
	method := ""
	if m := r.PostFormValue("remember"); m == "on" {
		method = "remember"
	}
	sts = ses.CreateSession(ctx, w, id, method)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.ByText(nil, config.StatusOK,
		"you just logged in, to return on main page click button below")
	api.Message(w, sts)
}

func (hu *hUser) Logout(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request,
) {
	if r.Method != "GET" {
		api.Message(w, object.ByCode(config.Code405))
		return
	}
	sts := ses.EndSession(w)
	if sts != nil {
		api.Message(w, sts)
		return
	}
	// w status
	sts = object.ByText(nil, config.StatusOK,
		"you just logged out, to return on main page click button below")
	api.Message(w, sts)
}
