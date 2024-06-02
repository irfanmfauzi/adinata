package repository

import (
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"context"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	baseRepo
}

func NewUserRepo(db *sqlx.DB) userRepo {
	return userRepo{
		baseRepo: baseRepo{db: db},
	}
}

func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (user entity.User, err error) {
	query := "SELECT id,email,password,role FROM users where email = $1"

	err = u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepo) InsertUser(ctx context.Context, tx TxProvider, req request.RegisterRequest) (int64, error) {
	user_id := int64(0)

	query := "INSERT INTO users (email, password, role) VALUES ($1,$2,$3) RETURNING id"

	err := u.DB(tx).GetContext(ctx, &user_id, query, req.Email, req.Password, req.Role)
	if err != nil {
		return user_id, err
	}

	return user_id, nil

}
