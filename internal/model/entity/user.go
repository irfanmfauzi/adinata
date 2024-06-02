package entity

type User struct {
	Id       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password,omitempty" json:"-"`
	Role     string `db:"role" json:"role"`
}
