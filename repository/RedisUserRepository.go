package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"magmar/model"

	"github.com/go-redis/redis/v8"
)

type redisUserRepository struct {
	client           *redis.Client
	userReadOnlyRepo UserReadOnlyRepository
}

// NewRedisUserRepository ...
func NewRedisUserRepository(client *redis.Client, userReadOnlyRepo UserReadOnlyRepository) UserReadOnlyRepository {
	return &redisUserRepository{
		client:           client,
		userReadOnlyRepo: userReadOnlyRepo,
	}
}

// GetUser ...
func (r *redisUserRepository) GetUser(ctx context.Context, uid string) (ruser *model.User, err error) {
	zlog.With(ctx).Infow("[New RedisRepository Request]", "uid", uid)
	userJSON, err := r.client.Get(ctx, fmt.Sprintf("%s:users:%s", redisPrefix, uid)).Bytes()
	if err == redis.Nil {
		zlog.With(ctx).Infow("GetUser Not Found", "uid", uid)
		if ruser, err = r.userReadOnlyRepo.GetUser(ctx, uid); err == nil {
			r.newUserToRedis(ctx, ruser)
		}
		return ruser, err
	} else if err != nil {
		zlog.With(ctx).Infow("GetUser Error", "uid", uid, "err", err)
		return nil, err
	}

	if err = json.Unmarshal(userJSON, &ruser); err != nil {
		return nil, err
	}
	return ruser, nil
}

// GetUserByEmail ...
func (r *redisUserRepository) GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error) {
	return r.userReadOnlyRepo.GetUserByEmail(ctx, email)
}

func (r *redisUserRepository) newUserToRedis(ctx context.Context, user *model.User) (err error) {
	zlog.With(ctx).Infow("[New RedisRepository Request]", "user", user)
	userJSON, err := json.Marshal(user)
	if err != nil {
		zlog.With(ctx).Errorw("newUserToRedis Error", "err", err)
		return err
	}

	if err = r.client.Set(ctx, fmt.Sprintf("%s:users:%s", redisPrefix, user.UID), userJSON, 0).Err(); err != nil {
		zlog.With(ctx).Errorw("newUserToRedis Error", "err", err)
	}
	return err
}
