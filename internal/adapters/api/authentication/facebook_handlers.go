package authentication

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

func (ha *hAuth) loginFacebook(w http.ResponseWriter, r *http.Request) {
	if ha.auth.Facebook.Empty {
		api.Message(w, object.ByText(errors.New("facebook authentication settings is null"), config.NotWorking,
			"facebook authentication"))
		return
	}
	sc := social{
		clientID: ha.auth.Facebook.ClientID,
		authURL:  config.FacebookAuthURL,
		redirect: ha.auth.Facebook.Redirect,
		scope:    "public_profile email",
	}
	ha.login(w, r, sc)
}

func (ha *hAuth) loginFacebookCallback(ctx context.Context, ses api.Middleware,
	w http.ResponseWriter, r *http.Request,
) {
	sc := social{
		name:         config.KeyFacebook,
		clientID:     ha.auth.Facebook.ClientID,
		clientSecret: ha.auth.Facebook.ClientSecret,
		tokenURL:     config.FacebookTokenURL,
		userURL:      config.FacebookUserURL,
		redirect:     ha.auth.Facebook.Redirect,
	}
	ha.callback(ctx, ses, w, r, sc)
}

func getFacebookData(urlStr string) ([]byte, object.Status) {
	client := http.Client{
		Timeout: config.TimeLimit2s,
	}
	response, err := client.Get(urlStr)
	if err != nil {
		return nil, object.ByCodeAndLog(config.Code500,
			err, "auth: API: request failed")
	}
	// Read the response as a byte slice
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, object.ByCodeAndLog(config.Code500,
			err, "auth: API: response read failed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("auth: API: response body close error: %v\n", err)
		}
	}(response.Body)
	return body, nil
}
