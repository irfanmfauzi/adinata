package repository

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"context"

	"github.com/jmoiron/sqlx"
)

type tagRepo struct {
	baseRepo
}

func NewTagRepo(db *sqlx.DB) tagRepo {
	return tagRepo{
		baseRepo: baseRepo{
			db: db,
		},
	}
}

func (t *tagRepo) UpsertTag(ctx context.Context, tx TxProvider, req []request.TagDbRequest) ([]int64, error) {

	tagIds := []int64{}
	upsertTagQuery := "INSERT INTO tags (label) VALUES (:tags) ON CONFLICT (label) DO UPDATE set label=excluded.label RETURNING id"

	upsertTagQuery, args, err := sqlx.Named(upsertTagQuery, req)
	if err != nil {
		return tagIds, err
	}

	err = tx.SelectContext(ctx, &tagIds, sqlx.Rebind(sqlx.DOLLAR, upsertTagQuery), args...)
	if err != nil {
		return tagIds, err
	}
	return tagIds, nil
}

func (t *tagRepo) InsertTag(ctx context.Context, tx TxProvider, req request.TagDbRequest) error {

	upsertTagQuery := "INSERT INTO tags (label) VALUES ($1)"

	_, err := t.DB(tx).ExecContext(ctx, upsertTagQuery, req.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (t *tagRepo) UpdateTag(ctx context.Context, tx TxProvider, req request.TagDbRequest, tagId int64) error {

	query := "UPDATE tags set label=$1 where id = $2"

	_, err := t.DB(tx).ExecContext(ctx, query, req.Tags, tagId)
	if err != nil {
		return err
	}

	return nil
}

func (t *tagRepo) DeleteTag(ctx context.Context, tx TxProvider, tagId int64) error {
	query := "Delete tags where id = $2"

	_, err := t.DB(tx).ExecContext(ctx, query, tagId)
	if err != nil {
		return err
	}

	return nil

}

func (t *tagRepo) GetTag(ctx context.Context) ([]entity.Tag, error) {
	query := "SELECT * FROM tags"

	result := []entity.Tag{}

	err := t.db.SelectContext(ctx, &result, query)
	if err != nil {
		return result, err
	}

	return result, nil

}
