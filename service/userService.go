package service

import (
	"context"
	"magmar/model"
	"magmar/repository"

	"github.com/juju/errors"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserService ...
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// GetUser ...
func (u *userUsecase) GetUser(ctx context.Context, uid string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Service Request]", "uid", uid)
	if ruser, err = u.userRepo.GetUser(ctx, uid); err != nil {
		zlog.With(ctx).Errorw("UserRepo GetUser Failed", "uid", uid, "err", err)
		return nil, err
	}

	return ruser, nil
}

// GetUserByEmail ...
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Service Request]", "email", email)
	if ruser, err = u.userRepo.GetUserByEmail(ctx, email); err != nil {
		zlog.With(ctx).Errorw("UserRepo GetUserByEmail Failed", "email", email, "err", err)
		return nil, err
	}

	return ruser, nil
}

// UpdateUser ...
func (u *userUsecase) UpdateUser(ctx context.Context, uid string, user *model.User) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Service Request]", "user", user)

	if !user.ValidateUpdateUser() {
		zlog.With(ctx).Warnw("ValidateUpdateUser failed", "user", user)
		return nil, errors.NotValidf("ValidateUpdateUser failed")
	}

	ruser, err = u.GetUser(ctx, uid)
	if err != nil {
		zlog.With(ctx).Errorw("UserRepo UpdateUser Failed", "err", err)
		return nil, err
	}

	ruser.UpdateUser(user)
	return u.userRepo.UpdateUser(ctx, ruser)
}

// DeleteUser ...
func (u *userUsecase) DeleteUser(ctx context.Context, uid string) (err error) {
	zlog.With(ctx).Infow("[New Service Request]", "uid", uid)

	return u.userRepo.DeleteUser(ctx, uid)
}
