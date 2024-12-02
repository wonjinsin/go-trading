package model

import (
	"gorm.io/gorm"
)

// DB ...
type DB struct {
	MainDB *gorm.DB
	ReadDB *gorm.DB
}

// WithMainDB ...
func (db *DB) WithMainDB() *gorm.DB {
	return db.MainDB
}

// WithReadDB ...
func (db *DB) WithReadDB() *gorm.DB {
	return db.ReadDB
}
