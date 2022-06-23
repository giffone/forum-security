package config

import (
	"time"
)

const (
	CookieSession        = "session"        // name for cookie
	CookieUserID         = "userID"         // name for cookie
	CookiePostID         = "postID"         // name for cookie
	LoginMinLength       = 2                // symbols
	PasswordMinLength    = 6                // symbols
	MaxPostsOnPage       = 15               // 15 post will show on main page
	TimeLimit20s         = 20 * time.Second
	TimeLimit10s         = 10 * time.Second
	TimeLimit5s          = 5 * time.Second
	TimeLimit2s          = 2 * time.Second
	ForumLayoutDate      = "January 2, 2006" // format for page
	MaxImageSize         = int64(20 << 20)   // 20Mb
	MaxImageSizeStr      = "20mb"
)
