package request

type PostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `db:"tags" json:"tags"`
}

type PostTagDbRequest struct {
	PostId int64 `db:"post_id"`
	TagId  int64 `db:"tag_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TagDbRequest struct {
	Tags string `db:"tags"`
}

type TagRequest struct {
	Label string `db:"label"`
}
