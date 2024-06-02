package repository

import (
	"adinata/internal/model/entity"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	baseRepo
}

func NewPostRepo(db *sqlx.DB) postRepo {
	return postRepo{
		baseRepo{
			db: db,
		},
	}
}

func (p *postRepo) InsertPost(ctx context.Context, tx TxProvider, title, content string, status entity.Status, publishTime sql.NullTime) (int64, error) {
	insertPostQuery := "INSERT INTO posts (title,content,status,publish_date) VALUES ($1, $2, $3, $4) RETURNING id"

	postId := int64(0)

	err := p.DB(tx).GetContext(ctx, &postId, insertPostQuery, title, content, string(status), publishTime)
	if err != nil {
		return postId, err
	}
	return postId, nil
}

func (p *postRepo) UpdatePost(ctx context.Context, tx TxProvider, postId int64, title, content string, status entity.Status, publishTime sql.NullTime) error {
	updatePostQuery := "UPDATE posts set title=$1,content=$2,status=$3,publish_date=$4 WHERE id = $5"

	_, err := tx.ExecContext(ctx, updatePostQuery, title, content, status, publishTime, postId)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepo) DeletePost(ctx context.Context, tx TxProvider, postId int64) error {
	deletePostQuery := "DELETE FROM posts where id = $1"

	_, err := tx.ExecContext(ctx, deletePostQuery, postId)
	if err != nil {
		return err
	}

	return nil
}

func (p *postRepo) GetPostById(ctx context.Context, postId int64) (entity.Post, error) {
	getPostByIdQuery := "select *, (select array_agg(label) from tags join post_tags pt on pt.tag_id = tags.id where pt.post_id  = $1) as tags from posts p where p.id = $1;"

	postData := entity.Post{}

	err := p.db.GetContext(ctx, &postData, getPostByIdQuery, postId)
	if err != nil {
		return postData, err
	}

	return postData, nil
}

func (p *postRepo) GetPosts(ctx context.Context, searchParams string) ([]entity.Post, error) {
	query := `
	select posts.*, array_agg(tags.label) as tags
	from
		posts
	join post_tags pt on pt.post_id = posts.id
	join tags on tags.id = pt.tag_id
	group by posts.id
	having
	    array_to_string(array_agg(tags.label), ';') ilike '%'||$1||'%';
	`
	result := []entity.Post{}
	err := p.db.SelectContext(ctx, &result, query, searchParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}
