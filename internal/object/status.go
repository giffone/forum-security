package object

import (
	"fmt"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/constant"
)

type Statuses struct {
	StatusBody string
	StatusCode int
	ReturnPage string
}

func NewStatuses() *Statuses {
	return new(Statuses)
}

type Status interface {
	Status() *Statuses
}

func (s *Statuses) Status() *Statuses {
	return s
}

func (s *Statuses) ByCodeAndLog(code int, err error, message string) {
	log.Printf("statusByCodeAndLog method: %s: %v\n", message, err)
	s.StatusBody = http.StatusText(code)
	s.StatusCode = code
}

func (s *Statuses) ByCode(code int) {
	s.StatusBody = http.StatusText(code)
	s.StatusCode = code
}

func (s *Statuses) ByText(err error, text string, args ...any) {
	sts := ByText(err, text, args)
	s.StatusBody = sts.StatusBody
	s.StatusCode = sts.StatusCode
}

func ByCodeAndLog(code int, err error, message string) *Statuses {
	log.Printf("%s: %v\n", message, err)
	return &Statuses{
		StatusBody: http.StatusText(code),
		StatusCode: code,
	}
}

func ByCode(code int) *Statuses {
	return &Statuses{
		StatusBody: http.StatusText(code),
		StatusCode: code,
	}
}

func ByText(err error, text string, args ...any) *Statuses {
	return byText(err, text, args)
}

func byText(err error, text string, args []any) *Statuses {
	sts := new(Statuses)
	if err != nil {
		e := err.Error()
		log.Printf("status by text: err: %s\n", e)
	}
	if len(args) == 0 {
		sts.StatusBody = text
	} else {
		sts.StatusBody = fmt.Sprintf(text, args...)
	}
	if text == constant.StatusOK {
		sts.StatusCode = constant.Code200
	} else if text == constant.StatusCreated {
		sts.StatusCode = constant.Code201
	} else {
		if sts.StatusCode == 0 {
			sts.StatusCode = constant.Code403
		}
	}
	return sts
}
