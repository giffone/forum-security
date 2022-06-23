package api

import (
	"context"
	"net/http"
	"time"

	"github.com/giffone/forum-security/internal/object"
)

type Handler interface {
	Register(ctx context.Context, router *http.ServeMux, session Middleware)
}

type Middleware interface { // middleware for handlers
	CreateSession(ctx context.Context, w http.ResponseWriter, id int, method string) object.Status
	Skip(ctx context.Context, fn func(context.Context, Middleware,
		http.ResponseWriter, *http.Request)) http.HandlerFunc
	CheckSession(ctx context.Context, fn func(context.Context,
		*object.CookieInfo, object.Status, http.ResponseWriter,
		*http.Request)) http.HandlerFunc
	EndSession(w http.ResponseWriter) object.Status
}

type Banned interface {
	Add(key string, banTill time.Time)
	UnlockBan(deadline time.Time, key string) (exist bool, expire bool)
}
