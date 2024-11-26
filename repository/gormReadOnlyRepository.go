package repository

import (
	"context"
	"magmar/model"

	"github.com/juju/errors"
	"gorm.io/gorm"
)

type gormUserReadOnlyRepository struct {
	conn *gorm.DB
}

// NewGormUserReadOnlyRepository ...
func NewGormUserReadOnlyRepository(conn *gorm.DB) UserRepository {
	return &gormUserRepository{conn: conn}
}

// GetUser ...
func (g *gormUserReadOnlyRepository) GetUser(ctx context.Context, uid string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Repository Request]", "uid", uid)

	scope := g.conn.WithContext(ctx)
	scope = scope.Where("users.uid = ?", uid).Find(&ruser)
	if err = scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetUser Error", "uid", uid, "err", err)
		return nil, err
	}
	if scope.RowsAffected == 0 {
		zlog.With(ctx).Errorw("GetUser Not Found", "uid", uid, "err", err)
		return nil, errors.UserNotFoundf("User is not exist")
	}

	return ruser, nil
}

// GetUserByEmail ...
func (g *gormUserReadOnlyRepository) GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Repository Request]", "email", email)

	scope := g.conn.WithContext(ctx)
	scope = scope.Where("users.email = ?", email).Find(&ruser)
	if err = scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetUserByEmail Error", "email", email, "err", err)
		return nil, err
	}
	if scope.RowsAffected == 0 {
		zlog.With(ctx).Errorw("GetUserByEmail Not Found", "email", email, "err", err)
		return nil, errors.UserNotFoundf("User is not exist")
	}

	return ruser, nil
}
