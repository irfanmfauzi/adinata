package service

import (
	"adinata/internal/repository"
	authService "adinata/internal/service/auth"
	postService "adinata/internal/service/post"
	tagService "adinata/internal/service/tag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	UserService authService.AuthServiceProvider
	PosService  postService.PostServiceProvider
	TagService  tagService.TagServiceProvider
}

var (
	database        = os.Getenv("DB_DATABASE")
	password        = os.Getenv("DB_PASSWORD")
	username        = os.Getenv("DB_USERNAME")
	port            = os.Getenv("DB_PORT")
	host            = os.Getenv("DB_HOST")
	serviceInstance Service
)

func New() Service {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepo(db)
	postRepo := repository.NewPostRepo(db)
	tagRepo := repository.NewTagRepo(db)
	postTagRepo := repository.NewPostTagRepo(db)

	userService := authService.NewAuthService(authService.AuthServiceConfig{
		UserRepo: &userRepo,
		Db:       db,
	})

	postService := postService.NewAuthService(postService.PostServiceConfig{
		Db:          db,
		PostRepo:    &postRepo,
		PostTagRepo: &postTagRepo,
		TagRepo:     &tagRepo,
	})

	tagService := tagService.NewTagService(tagService.TagServiceConfig{
		Db:      db,
		TagRepo: &tagRepo,
	})

	return Service{
		UserService: &userService,
		PosService:  &postService,
		TagService:  &tagService,
	}
}
