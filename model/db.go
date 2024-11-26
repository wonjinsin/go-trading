package model

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// DB ...
type DB struct {
	MainDB *gorm.DB
	ReadDB *gorm.DB
	Redis  *redis.Client
}

// WithMainDB ...
func (db *DB) WithMainDB() *gorm.DB {
	return db.MainDB
}

// WithReadDB ...
func (db *DB) WithReadDB() *gorm.DB {
	return db.ReadDB
}

// WithRedis ...
func (db *DB) WithRedis() *redis.Client {
	return db.Redis
}
