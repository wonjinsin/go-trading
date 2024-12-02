package service

import (
	"context"
	"log"
	"magmar/config"
	"magmar/repository"
	"magmar/util"
	"os"

	"magmar/model"
)

var zlog *util.Logger

func init() {
	var err error
	zlog, err = util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[service] err[%s]", err.Error())
		os.Exit(1)
	}
}

// Service ...
type Service struct {
	User UserService
	Auth AuthService
}

// Init ...
func Init(conf *config.ViperConfig, repo *repository.Repository) (*Service, error) {
	userSvc := NewUserService(repo.User)
	authSvc := NewAuthService(conf, repo.User)
	return &Service{
		User: userSvc,
		Auth: authSvc,
	}, nil
}

// UserService ...
type UserService interface {
	GetUser(ctx context.Context, id string) (ruser *model.User, err error)
	GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error)
	UpdateUser(ctx context.Context, uid string, user *model.User) (ruser *model.User, err error)
	DeleteUser(ctx context.Context, id string) (err error)
}

// AuthService ...
type AuthService interface {
	Signup(ctx context.Context, signup *model.Signup) (token *model.Token, err error)
	Signin(ctx context.Context, signin *model.Signin) (token *model.Token, err error)
}
