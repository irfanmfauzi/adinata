package tagService

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"adinata/internal/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

type TagServiceProvider interface {
	CreateTag(ctx context.Context, req request.TagRequest) error
	GetTag(ctx context.Context) ([]entity.Tag, error)
	UpdateTag(ctx context.Context, req request.TagRequest, tagId int64) error
	DeleteTag(ctx context.Context, tagId int64) error
}

type TagServiceConfig struct {
	Db      *sqlx.DB
	TagRepo repository.TagRepoProvider
}

func NewTagService(cfg TagServiceConfig) tagService {
	return tagService{
		db:      cfg.Db,
		tagRepo: cfg.TagRepo,
	}
}

type tagService struct {
	db      *sqlx.DB
	tagRepo repository.TagRepoProvider
}

func (t *tagService) CreateTag(ctx context.Context, req request.TagRequest) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = t.tagRepo.InsertTag(ctx, tx, request.TagDbRequest{Tags: req.Label})
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (t *tagService) GetTag(ctx context.Context) ([]entity.Tag, error) {
	return t.tagRepo.GetTag(ctx)
}

func (t *tagService) UpdateTag(ctx context.Context, req request.TagRequest, tagId int64) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = t.tagRepo.UpdateTag(ctx, tx, request.TagDbRequest{Tags: req.Label}, tagId)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (t *tagService) DeleteTag(ctx context.Context, tagId int64) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = t.tagRepo.DeleteTag(ctx, tx, tagId)
	if err != nil {
		return err
	}

	return nil

}
