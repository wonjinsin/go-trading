package service

import (
	"context"
	"magmar/config"
	"magmar/model"
	"magmar/repository"

	"github.com/juju/errors"
)

type authUsecase struct {
	conf     *config.ViperConfig
	userRepo repository.UserRepository
}

// NewAuthService ...
func NewAuthService(conf *config.ViperConfig, userRepo repository.UserRepository) AuthService {
	return &authUsecase{
		conf:     conf,
		userRepo: userRepo,
	}
}

// Signup ...
func (a *authUsecase) Signup(ctx context.Context, signup *model.Signup) (token *model.Token, err error) {
	zlog.With(ctx).Infow("[New Service Request]", "signup", signup)
	if _, err := a.userRepo.GetUserByEmail(ctx, signup.Email); err == nil {
		zlog.With(ctx).Errorw("UserRepo UserExist", "signup", signup)
		return nil, errors.AlreadyExistsf("User")
	}

	user, err := model.NewUserBySignup(signup)
	if err != nil {
		zlog.With(ctx).Errorw("NewUserBySignup failed", "signup", signup)
		return nil, errors.NotValidf("Signup")
	}

	user, err = a.userRepo.NewUser(ctx, user)
	if err != nil {
		zlog.With(ctx).Warnw("Update by signup", "user", user, "signup", signup)
		return nil, err
	}
	return a.newToken(ctx, user)
}

// Signin ...
func (a *authUsecase) Signin(ctx context.Context, signin *model.Signin) (token *model.Token, err error) {
	zlog.With(ctx).Infow("[New Service Request]", "signin", signin)
	user, err := a.userRepo.GetUserByEmail(ctx, signin.Email)
	if err != nil {
		zlog.With(ctx).Warnw("GetUserByEmail failed", "signin", signin)
		return nil, err
	}
	return a.newToken(ctx, user)
}

func (a *authUsecase) newToken(ctx context.Context, user *model.User) (*model.Token, error) {
	tokenClaim := model.NewUserClaim(a.conf.GetString("projectName"), a.conf.GetString("domain"), user)
	accessToken, err := tokenClaim.GenerateToken(a.conf.Get("prvTokenKey").([]byte))
	if err != nil {
		zlog.With(ctx).Errorw("GenerateToken failed", "user", user)
		return nil, err
	}
	return model.NewToken(accessToken), nil
}
