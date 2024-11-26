package repository

import (
	"context"
	"magmar/model"
	"magmar/util"

	"github.com/juju/errors"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	conn *gorm.DB
}

// NewGormUserRepository ...
func NewGormUserRepository(conn *gorm.DB) UserRepository {
	migrations := []interface{}{
		&model.User{},
	}
	if err := conn.Set("gorm:table_options", util.DBCharsetOption).Migrator().AutoMigrate(migrations...); err != nil {
		zlog.Panicw("NewGormDealRepository Unable to AutoMigrate DealRepository", "err", err)
	}

	return &gormUserRepository{conn: conn}
}

// NewUser ...
func (g *gormUserRepository) NewUser(ctx context.Context, user *model.User) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Repository Request]", "user", user)

	scope := g.conn.WithContext(ctx)
	if err = scope.Create(&user).Error; err != nil {
		zlog.With(ctx).Errorw("NewUser Error", "err", err)
		return nil, err
	}

	return user, nil
}

// GetUser ...
func (g *gormUserRepository) GetUser(ctx context.Context, uid string) (ruser *model.User, err error) {
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
func (g *gormUserRepository) GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Repository Request]", "email", email)

	scope := g.conn.WithContext(ctx)
	scope = scope.Where("users.email = ?", email).Find(&ruser)
	if err = scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetUserByEmail Error", "email", email, "err", err)
		return nil, err
	}
	if scope.RowsAffected == 0 {
		zlog.With(ctx).Warnw("GetUserByEmail Not Found", "email", email, "err", err)
		return nil, errors.UserNotFoundf("User is not exist")
	}

	return ruser, nil
}

// UpdateUser ...
func (g *gormUserRepository) UpdateUser(ctx context.Context, user *model.User) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New Repository Service]", "user", user)

	scope := g.conn.WithContext(ctx)
	if err = scope.Updates(user).Error; err != nil {
		zlog.With(ctx).Errorw("UpdateUser Error", "user", user, "err", err)
		return nil, err
	}

	return user, nil
}

// DeleteUser ...
func (g *gormUserRepository) DeleteUser(ctx context.Context, uid string) (err error) {
	zlog.With(ctx).Infow("[New Repository Request]", "uid", uid)

	scope := g.conn.WithContext(ctx)
	if err = scope.Where("uid = ?", uid).Delete(&model.User{}).Error; err != nil {
		zlog.With(ctx).Errorw("DeleteUser Error", "uid", uid, "err", err)
		return err
	}

	return nil
}
