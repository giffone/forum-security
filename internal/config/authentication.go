package config

import (
	"fmt"
	"log"
	"os"

	"github.com/giffone/forum-security/pkg/password"
	"github.com/giffone/forum-security/pkg/read_env"
)

const (
	AuthenticationPath = "secure/authentication/.env"

	GithubAuthURL  = "https://github.com/login/oauth/authorize"
	GithubTokenURL = "https://github.com/login/oauth/access_token"
	GithubUserURL  = "https://api.github.com/user"

	/*------------------------------------------------------*/

	FacebookAuthURL  = "https://www.facebook.com/dialog/oauth"
	FacebookTokenURL = "https://graph.facebook.com/oauth/access_token"
	FacebookUserURL  = "https://graph.facebook.com/me?fields=id,name,email&access_token="

	/*------------------------------------------------------*/

	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleUserURL  = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
)

type Auth struct {
	Github     *Github
	Facebook   *Facebook
	Google     *Google
	StateToken string // unique key for safe redirect url
}

type Github struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

type Facebook struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

type Google struct {
	ClientID     string
	ClientSecret string
	Redirect     string
	Empty        bool
}

func NewSocialToken(home string) *Auth {
	_ = read_env.ReadEnv(AuthenticationPath)
	auth := &Auth{
		StateToken: password.Generate(),
		Github: &Github{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, URLLoginGithubCallback),
		},
		Facebook: &Facebook{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, URLLoginFacebookCallback),
		},
		Google: &Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Redirect:     fmt.Sprintf("%s%s", home, URLLoginGoogleCallback),
		},
	}
	os.Clearenv()
	if auth.Github.ClientID == "" {
		log.Println("Missing Github Client ID")
		auth.Github.Empty = true
	}
	if auth.Github.ClientSecret == "" {
		log.Println("Missing Github Client Secret")
		auth.Github.Empty = true
	}
	if auth.Facebook.ClientID == "" {
		log.Println("Missing Facebook Client ID")
		auth.Facebook.Empty = true
	}
	if auth.Facebook.ClientSecret == "" {
		log.Println("Missing Facebook Client Secret")
		auth.Facebook.Empty = true
	}
	if auth.Google.ClientID == "" {
		log.Println("Missing Google Client ID")
		auth.Google.Empty = true
	}
	if auth.Google.ClientSecret == "" {
		log.Println("Missing Google Client Secret")
		auth.Google.Empty = true
	}
	return auth
}
