package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
)

func (ha *hAuth) loginGoogle(w http.ResponseWriter, r *http.Request) {
	if ha.auth.Github.Empty {
		api.Message(w, object.ByText(errors.New("github authentication settings is null"), constant.NotWorking,
			"github authentication"))
		return
	}
	sc := social{
		clientID: ha.auth.Google.ClientID,
		authURL:  constant.GoogleAuthURL,
		redirect: ha.auth.Google.Redirect,
		scope:    "profile email",
	}
	ha.login(w, r, sc)
}

func (ha *hAuth) loginGoogleCallback(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request,
) {
	sc := social{
		name:         constant.KeyGoogle,
		clientID:     ha.auth.Google.ClientID,
		clientSecret: ha.auth.Google.ClientSecret,
		tokenURL:     constant.GoogleTokenURL,
		userURL:      constant.GoogleUserURL,
		redirect:     ha.auth.Google.Redirect,
	}
	ha.callback(ctx, ses, w, r, sc)
}
