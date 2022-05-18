package filemaker

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/service"
)

type sFileMaker struct {
	repo repository.Repo
}

func NewService(repo repository.Repo) service.FileMaker {
	return &sFileMaker{
		repo: repo,
	}
}

func (fm *sFileMaker) CreateFile(ctx context.Context, d *dto.FileMaker) object.Status {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()

	file, err := os.Create(d.Path)
	if err != nil {
		return object.ByCodeAndLog(constant.Code500,
			err, "create file")
	}
	_, err = io.Copy(file, bytes.NewReader(d.Src.Body))
	if cerr := file.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(file.Name())
		return object.ByCodeAndLog(constant.Code500,
			err, "create file: read or close")
	}
	// for web need to cut path
	a := strings.TrimPrefix(d.Path, "internal/web")
	d.Path = a
	// create image path in db
	_, sts := fm.repo.Create(ctx2, d)
	if sts != nil {
		return sts
	}
	return nil
}
