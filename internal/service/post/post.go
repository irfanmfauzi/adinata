package postService

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"adinata/internal/repository"
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostServiceProvider interface {
	CreatePost(ctx context.Context, req request.PostRequest) error

	UpdatePost(ctx context.Context, req request.PostRequest, postId int64) error

	DeletePost(ctx context.Context, posId int64) error

	GetPostDetail(ctx context.Context, posId int64) (entity.Post, error)
}

type PostServiceConfig struct {
	Db          *sqlx.DB
	PostRepo    repository.PostRepoProvider
	PostTagRepo repository.PostTagRepoProvider
	TagRepo     repository.TagRepoProvider
}

type postService struct {
	db          *sqlx.DB
	postRepo    repository.PostRepoProvider
	postTagRepo repository.PostTagRepoProvider
	tagRepo     repository.TagRepoProvider
}

func NewAuthService(cfg PostServiceConfig) postService {
	return postService{
		db:          cfg.Db,
		postRepo:    cfg.PostRepo,
		postTagRepo: cfg.PostTagRepo,
		tagRepo:     cfg.TagRepo,
	}
}

func (p *postService) CreatePost(ctx context.Context, req request.PostRequest) error {
	user := ctx.Value("user").(map[string]interface{})

	role := user["role"].(string)

	status := entity.Draft
	publishTime := sql.NullTime{Time: time.Now().UTC()}

	if role == "admin" {
		status = entity.Publish
		publishTime.Valid = true
	}

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tagVals := []request.TagDbRequest{}

	for _, v := range req.Tags {
		tagVals = append(tagVals, request.TagDbRequest{Tags: v})
	}

	tagIds, err := p.tagRepo.UpsertTag(ctx, tx, tagVals)
	if err != nil {
		return err
	}

	postId, err := p.postRepo.InsertPost(ctx, tx, req.Title, req.Content, status, publishTime)
	if err != nil {
		return err
	}

	reqPostTag := make([]request.PostTagDbRequest, len(tagIds))

	for i, v := range tagIds {
		reqPostTag[i] = request.PostTagDbRequest{
			PostId: postId,
			TagId:  v,
		}
	}

	err = p.postTagRepo.InsertPostTag(ctx, tx, reqPostTag)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (p *postService) UpdatePost(ctx context.Context, req request.PostRequest, postId int64) error {
	user := ctx.Value("user").(map[string]interface{})
	role := user["role"].(string)

	status := entity.Draft
	publishTime := sql.NullTime{Time: time.Now().UTC()}

	if role == "admin" {
		status = entity.Publish
		publishTime.Valid = true
	}

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.postRepo.UpdatePost(ctx, tx, postId, req.Title, req.Content, status, publishTime)
	if err != nil {
		slog.Error("Error Update Posts ", "Error", err)
		return err
	}

	tagVals := []request.TagDbRequest{}

	mapTags := make(map[string]struct{})

	for _, v := range req.Tags {
		tagVals = append(tagVals, request.TagDbRequest{Tags: v})
		mapTags[v] = struct{}{}
	}

	tagIds, err := p.tagRepo.UpsertTag(ctx, tx, tagVals)
	if err != nil {
		return err
	}

	postTagVals := make([]request.PostTagDbRequest, len(tagIds))
	for i, v := range tagIds {
		postTagVals[i] = request.PostTagDbRequest{PostId: postId, TagId: v}
	}

	_, err = p.postTagRepo.UpsertPostTag(ctx, tx, postTagVals)
	if err != nil {
		return err
	}

	postTags, err := p.postTagRepo.GetPostTagByPostId(ctx, tx, postId)
	if err != nil {
		return err
	}

	listDeletedTagId := []int64{}

	for _, v := range postTags {
		_, ok := mapTags[v.Label]
		if ok {
			delete(mapTags, v.Label)
		} else {
			listDeletedTagId = append(listDeletedTagId, v.TagId)
		}
	}

	if len(listDeletedTagId) > 0 {
		err = p.postTagRepo.DeletePostTagByPostIdAndTagIds(ctx, tx, postId, listDeletedTagId)
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil

}

func (p *postService) DeletePost(ctx context.Context, posId int64) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to Begin Transaction", "Error", err)
		return err
	}
	defer tx.Rollback()

	err = p.postRepo.DeletePost(ctx, tx, posId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to Commit", "Error", err)
		return err
	}
	return nil
}

func (p *postService) GetPostDetail(ctx context.Context, postId int64) (entity.Post, error) {
	return p.postRepo.GetPostById(ctx, postId)
}
