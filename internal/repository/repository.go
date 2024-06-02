package repository

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"context"
	"database/sql"
)

type UserRepoProvider interface {
	FindUserByEmail(ctx context.Context, email string) (user entity.User, err error)
	InsertUser(ctx context.Context, tx TxProvider, req request.RegisterRequest) (int64, error)
}

type TagRepoProvider interface {
	UpsertTag(ctx context.Context, tx TxProvider, req []request.TagDbRequest) ([]int64, error)
	InsertTag(ctx context.Context, tx TxProvider, req request.TagDbRequest) error
	UpdateTag(ctx context.Context, tx TxProvider, req request.TagDbRequest, tagId int64) error
	DeleteTag(ctx context.Context, tx TxProvider, tagId int64) error
	GetTag(ctx context.Context) ([]entity.Tag, error)
}

type PostRepoProvider interface {
	InsertPost(ctx context.Context, tx TxProvider, title, content string, status entity.Status, publishTime sql.NullTime) (int64, error)
	UpdatePost(ctx context.Context, tx TxProvider, postId int64, title, content string, status entity.Status, publishTime sql.NullTime) error
	DeletePost(ctx context.Context, tx TxProvider, postId int64) error
	GetPostById(ctx context.Context, postId int64) (entity.Post, error)
	GetPosts(ctx context.Context, searchParams string) ([]entity.Post, error)
}

type PostTagRepoProvider interface {
	GetPostTagByPostId(ctx context.Context, tx TxProvider, postId int64) ([]entity.PostTag, error)
	InsertPostTag(ctx context.Context, tx TxProvider, req []request.PostTagDbRequest) error
	UpsertPostTag(ctx context.Context, tx TxProvider, req []request.PostTagDbRequest) ([]int64, error)
	DeletePostTagByPostIdAndTagIds(ctx context.Context, tx TxProvider, postId int64, tagIds []int64) error
}

type TxProvider interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type QueryProvider interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
