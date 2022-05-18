package category

import (
	"context"
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/object/model"
	"github.com/giffone/forum-security/internal/service"
)

type sCategory struct {
	repo repository.Repo
}

func NewService(repo repository.Repo) service.Category {
	return &sCategory{repo: repo}
}

func (sc *sCategory) Create(ctx context.Context, dto *dto.Category) object.Status {
	return nil
}

func (sc *sCategory) Delete(ctx context.Context, id int) object.Status {
	return nil
}

func (sc *sCategory) GetList(ctx context.Context, m model.Models) (interface{}, object.Status) {
	ctx2, cancel := context.WithTimeout(ctx, constant.TimeLimitDB)
	defer cancel()
	err := sc.repo.GetList(ctx2, m)
	if err != nil {
		return nil, err
	}
	categories := m.Return().Categories
	if len(categories.Slice) == 0 {
		return categories.IfNil(), nil
	}
	return categories.Slice, nil
}

func (sc *sCategory) GetFor(ctx context.Context, pc model.PostOrComment) object.Status {
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		categories := model.NewCategories(nil, nil)
		categories.MakeKeys(constant.KeyPost, id) // key - post, only post have category
		sts := sc.repo.GetList(ctx, categories)
		if sts != nil {
			return sts
		}
		if len(categories.Slice) == 0 {
			pc.Add(constant.KeyCategory, i, categories.IfNil())
		} else {
			pc.Add(constant.KeyCategory, i, categories.Slice)
		}
	}
	return nil
}

func (sc *sCategory) GetForChan(ctx context.Context, pc model.PostOrComment, channel chan object.Status) {
	for i := 0; i < pc.LSlice(); i++ {
		id := pc.PostOrCommentID(i)
		categories := model.NewCategories(nil, nil)
		categories.MakeKeys(constant.KeyPost, id) // key - post, only post have category
		sts := sc.repo.GetList(ctx, categories)
		if sts != nil {
			channel <- sts
		}
		if len(categories.Slice) == 0 {
			pc.Add(constant.KeyCategory, i, categories.IfNil())
		} else {
			pc.Add(constant.KeyCategory, i, categories.Slice)
		}
	}
	channel <- nil
}
