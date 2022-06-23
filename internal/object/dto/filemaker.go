package dto

import (
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
)

type FileMaker struct {
	Path string
	Type string
	Src  *FileSrc
	Obj  object.Obj
}

type FileSrc struct {
	Body []byte
	MIME string
}

func NewFileMaker(st *object.Settings, sts *object.Statuses, ck *object.CookieInfo) *FileMaker {
	fm := &FileMaker{
		Src: &FileSrc{},
	}
	fm.Obj.NewObjects(st, sts, ck)
	return fm
}

func (fm *FileMaker) MakeKeys(key string, data ...any) {
	if key != "" {
		fm.Obj.St.Key[key] = data
	} else {
		fm.Obj.St.Key[config.KeyPost] = []any{0}
	}
}

func (fm *FileMaker) Create() *object.QuerySettings {
	qs := new(object.QuerySettings)
	qs.QueryName = config.QueInsert4
	if value, ok := fm.Obj.St.Key[config.KeyPost]; ok {
		qs.QueryFields = []interface{}{
			config.TabFiles,
			config.FieldIdVariety,
			config.FieldVariety,
			config.FieldPath,
			config.FieldMIME,
		}
		qs.Fields = value
		qs.Fields = append(qs.Fields,
			config.KeyPost,
			fm.Path,
			fm.Type)
	}
	return qs
}

func (i *FileMaker) Delete() *object.QuerySettings {
	return &object.QuerySettings{}
}
