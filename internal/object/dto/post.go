package dto

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
)

type Post struct {
	Title      string
	Body       string
	Image      *FileMaker
	Categories *CategoryPost
	Obj        object.Obj
}

func NewPost(st *object.Settings, sts *object.Statuses, ck *object.Cookie) *Post {
	p := new(Post)
	p.Categories = NewCategoryPost()
	p.Obj.NewObjects(st, sts, ck)
	return p
}

func (p *Post) Add(r *http.Request) bool {
	// get user id
	sts := p.Obj.Ck.CookieUserIDRead(r)
	if sts != nil {
		p.Obj.Sts = sts.Status()
		return false
	}
	reader, err := r.MultipartReader()
	if err != nil {
		p.Obj.Sts.ByCodeAndLog(constant.Code400,
			err, "dto: create post: multipartReader")
		return false
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.Obj.Sts.ByCodeAndLog(constant.Code500,
				err, "dto: create post: multipartReader")
			return false
		}
		if !p.add(part) {
			return false
		}
	}
	return true
}

func (p *Post) add(part *multipart.Part) bool {
	pr := part.FormName()
	log.Println("pr is ", pr)
	switch pr {
	case "":
		return true // just skip and continue "for loop"
	case "title":
		text := make([]byte, 512)
		n, err := part.Read(text)
		if err != nil && err != io.EOF {
			p.Obj.Sts.ByCodeAndLog(constant.Code500,
				err, "dto: create post: form read")
			return false
		}
		p.Title = string(text[:n])
	case "body text":
		text := make([]byte, 5<<20)
		n, err := part.Read(text)
		if err != nil && err != io.EOF {
			p.Obj.Sts.ByCodeAndLog(constant.Code500,
				err, "dto: create post: form read")
			return false
		}
		p.Body = string(text[:n])
	case "categories":
		text := make([]byte, 512)
		n, err := part.Read(text)
		if err != nil && err != io.EOF {
			p.Obj.Sts.ByCodeAndLog(constant.Code500,
				err, "dto: create post: form read")
			return false
		}
		p.Categories.Slice = append(p.Categories.Slice, string(text[:n]))
	case "body image":
		name := part.FileName()
		if name == "" {
			return true // just skip and continue "for loop"
		}
		p.Image = NewFileMaker(nil, nil, nil)
		var b bytes.Buffer
		_, err := io.Copy(&b, part)
		if err != nil && err != io.EOF {
			p.Obj.Sts.ByCodeAndLog(constant.Code400,
				err, "dto: image: io copy")
			return false
		}
		p.Image.Src.MIME = part.Header.Get("Content-Type")
		p.Image.Src.Body = b.Bytes()
	default:
		p.Obj.Sts.ByCodeAndLog(constant.Code400,
			errors.New(pr), "unexpected form name")
		return false
	}
	return true
}

func (p *Post) Valid() bool {
	log.Println("in valid")
	if p.Obj.Sts.StatusBody != "" {
		return false
	}
	// delete space for check an any symbol
	body := strings.TrimSpace(p.Body)
	if body == "" {
		p.Obj.Sts.ByText(nil, constant.TooShort, "post", "one")
		return false
	}
	body = strings.TrimSpace(p.Title)
	if body == "" {
		if len(p.Body) > 20 {
			p.Title = p.Body[0:19] + "..."
		} else {
			p.Title = p.Body
		}
	}
	if p.Image != nil {
		// check for type file
		if p.Image.Src.MIME != "image/gif" && p.Image.Src.MIME != "image/png" && p.Image.Src.MIME != "image/jpeg" {
			p.Obj.Sts.ByText(nil, constant.ImageNotAllowed, p.Image.Src.MIME)
			p.Obj.Sts.StatusCode = constant.Code400
			return false
		} else {
			p.Image.Type = "image"
		}
		// check for limit size
		if int64(len(p.Image.Src.Body)) > constant.MaxImageSize {
			p.Obj.Sts.ByText(nil, constant.FileSizeToBig, constant.MaxImageSizeStr)
			p.Obj.Sts.StatusCode = constant.Code400
			return false
		}

	}
	return true
}

func (p *Post) Create() *object.QuerySettings {
	return &object.QuerySettings{
		QueryName: constant.QueInsert4,
		QueryFields: []interface{}{
			constant.TabPosts,
			constant.FieldUser,
			constant.FieldTitle,
			constant.FieldBody,
			constant.FieldCreated,
		},
		Fields: []interface{}{
			p.Obj.Ck.User,
			p.Title,
			p.Body,
			time.Now(),
		},
	}
}

func (p *Post) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}
