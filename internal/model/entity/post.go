package entity

import (
	"database/sql"

	"github.com/lib/pq"
)

type Status string

const (
	Draft   Status = "draft"
	Publish Status = "publish"
)

type Post struct {
	Id          int64          `db:"id" json:"id"`
	Title       string         `db:"title" json:"title"`
	Content     string         `db:"content" json:"content"`
	Tags        pq.StringArray `db:"tags" json:"tags"`
	PublishDate sql.NullTime   `db:"publish_date" json:"publish_date"`
	Status      Status         `db:"status" json:"status"`
}

type Tag struct {
	Id    int64   `db:"id"`
	Label string  `db:"label"`
	Posts []int64 `db:"posts"`
}

type PostTag struct {
	Id     int64  `db:"id"`
	PostId int64  `db:"post_id"`
	TagId  int64  `db:"tag_id"`
	Label  string `db:"label"`
}
