package repository

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"context"

	"github.com/jmoiron/sqlx"
)

type postTagRepo struct {
	baseRepo
}

func NewPostTagRepo(db *sqlx.DB) postTagRepo {
	return postTagRepo{
		baseRepo{
			db: db,
		},
	}
}

func (p *postTagRepo) InsertPostTag(ctx context.Context, tx TxProvider, req []request.PostTagDbRequest) error {
	insertPostTagQuery := "INSERT INTO post_tags (post_id,tag_id) VALUES (:post_id,:tag_id)"

	insertPostTagQuery, argsPostTag, err := sqlx.Named(insertPostTagQuery, req)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sqlx.Rebind(sqlx.DOLLAR, insertPostTagQuery), argsPostTag...)
	if err != nil {
		return err
	}
	return nil
}

func (p *postTagRepo) UpsertPostTag(ctx context.Context, tx TxProvider, req []request.PostTagDbRequest) ([]int64, error) {
	upsertPostTag := "INSERT INTO post_tags (post_id,tag_id) VALUES (:post_id,:tag_id) ON CONFLICT (post_id,tag_id) DO NOTHING RETURNING id"
	postTagId := []int64{}

	upsertPostTag, argsPostTag, err := sqlx.Named(upsertPostTag, req)
	if err != nil {
		return postTagId, err
	}

	err = tx.SelectContext(ctx, &postTagId, sqlx.Rebind(sqlx.DOLLAR, upsertPostTag), argsPostTag...)
	if err != nil {
		return postTagId, err
	}

	return postTagId, nil
}

func (p *postTagRepo) DeletePostTagByPostIdAndTagIds(ctx context.Context, tx TxProvider, postId int64, tagIds []int64) error {
	deletePostTagQuery := "DELETE FROM post_tags where post_id = $1 and tag_id = ANY($2)"
	_, err := tx.ExecContext(ctx, deletePostTagQuery, postId, tagIds)
	if err != nil {
		return err
	}
	return nil
}

func (p *postTagRepo) GetPostTagByPostId(ctx context.Context, tx TxProvider, postId int64) ([]entity.PostTag, error) {
	query := "SELECT post_tags.id,post_id,tag_id, label FROM post_tags join tags on tags.id = post_tags.tag_id where post_id = $1"
	result := []entity.PostTag{}

	err := p.DB(tx).SelectContext(ctx, &result, query, postId)
	if err != nil {
		return result, err
	}

	return result, nil
}
