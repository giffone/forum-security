package api

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type ParseExecute struct {
	PathTmpl   []string
	DefineTmpl string
	Data       map[string]interface{}
	tmpl       *template.Template
}

func NewParseExecute(define string) *ParseExecute {
	switch define {
	case "index":
		return &ParseExecute{
			PathTmpl: []string{
				config.PathIndex, config.PathHeaderObj,
				config.PathFooterObj, config.PathPostObj,
				config.PathCategoriesObj, config.PathLikesObj,
				config.PathPostObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "post":
		return &ParseExecute{
			PathTmpl: []string{
				config.PathPost, config.PathHeaderObj,
				config.PathPostObj, config.PathCategoriesObj,
				config.PathCommentsObj, config.PathLikesObj,
				config.PathFooterObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "account":
		return &ParseExecute{
			PathTmpl: []string{
				config.PathAccount, config.PathAccountUser,
				config.PathHeaderObj, config.PathPostObj,
				config.PathCategoriesObj, config.PathCommentsObj,
				config.PathLikesObj, config.PathFooterObj,
			},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "login":
		return &ParseExecute{
			PathTmpl:   []string{config.PathLoginObj},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	case "message":
		return &ParseExecute{
			PathTmpl:   []string{config.PathMessage},
			DefineTmpl: define,
			Data:       make(map[string]interface{}),
		}
	default:
		return &ParseExecute{
			DefineTmpl: "nil",
		}
	}
}

func (pe *ParseExecute) Parse() (*ParseExecute, object.Status) {
	if pe.DefineTmpl == "nil" {
		return nil, object.ByCodeAndLog(config.Code500,
			nil, "parseFile: unknown define")
	}
	var err error
	myFunc := template.FuncMap{
		"dateForum": func(t time.Time) string {
			return t.Format(config.ForumLayoutDate)
		},
	}
	pe.tmpl, err = template.New("").Funcs(myFunc).ParseFiles(pe.PathTmpl...)
	if err != nil {
		return nil, object.ByCodeAndLog(config.Code500,
			err, "parseFile:")
	}
	return pe, nil
}

func (pe *ParseExecute) Execute(w http.ResponseWriter, code int) {
	if code != 0 {
		w.WriteHeader(code)
		log.Printf("execute: writing code to w: %d", code)
	}

	err := pe.tmpl.ExecuteTemplate(w, pe.DefineTmpl, pe.Data)
	// log.Printf("execute: w header is: %v", w)
	if err != nil {
		sts := object.ByCodeAndLog(config.Code500,
			err, "executeTemplate:")
		Message(w, sts)
	}
}
