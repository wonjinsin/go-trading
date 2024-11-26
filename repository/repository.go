package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"magmar/config"
	"magmar/model"
	"magmar/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var zlog *util.Logger
var redisPrefix string

type dbLogger struct {
	*util.Logger
}

func (dl *dbLogger) LogMode(l logger.LogLevel) logger.Interface {
	return dl
}

func (dl *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Info(msg, data)
}

func (dl *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Warn(msg, data)
}

func (dl *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Error(msg, data)
}

func (dl *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		dl.Logger.With(ctx).Infow(err.Error(), "elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
	} else {
		dl.Logger.With(ctx).Infow("", "elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
	}
}

func init() {
	var err error
	zlog, err = util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[service] err[%s]", err.Error())
		os.Exit(1)
	}
}

// Repository ...
type Repository struct {
	User         UserRepository
	UserReadOnly UserReadOnlyRepository
}

// RedisRepository ...
type RedisRepository struct {
	UserReadOnly UserReadOnlyRepository
}

// Init ...
func Init(magmar *config.ViperConfig) (*Repository, *RedisRepository, error) {
	mysqlConn, err := mysqlConnect(magmar, "database")
	if err != nil {
		return nil, nil, err
	}

	mysqlReadOnlyConn, err := mysqlConnect(magmar, "readOnlyDatabase")
	if err != nil {
		return nil, nil, err
	}

	redisPrefix = magmar.GetString("projectName")
	redisConn, err := util.RedisConnect(magmar, zlog)
	if err != nil {
		return nil, nil, err
	}

	db := &model.DB{
		MainDB: mysqlConn,
		ReadDB: mysqlReadOnlyConn,
		Redis:  redisConn,
	}

	userRepo := NewGormUserRepository(db.MainDB)
	userReadOnlyRepo := NewGormUserReadOnlyRepository(db.ReadDB)

	redisUserRepo := NewRedisUserRepository(db.Redis, userRepo)

	return &Repository{
			User:         userRepo,
			UserReadOnly: userReadOnlyRepo,
		}, &RedisRepository{
			UserReadOnly: redisUserRepo,
		}, nil
}

func mysqlConnect(magmar *config.ViperConfig, prefix string) (mysql *gorm.DB, err error) {
	return gorm.Open(getDialector(magmar, prefix), &gorm.Config{})
}

func getDialector(magmar *config.ViperConfig, prefix string) gorm.Dialector {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=UTC",
		magmar.GetString(fmt.Sprintf("%s.username", prefix)),
		magmar.GetString(fmt.Sprintf("%s.password", prefix)),
		magmar.GetString(fmt.Sprintf("%s.host", prefix)),
		magmar.GetInt(fmt.Sprintf("%s.port", prefix)),
		magmar.GetString(fmt.Sprintf("%s.dbname", prefix)),
	)

	return mysql.Open(dbURI)
}

func getConfig(magmar *config.ViperConfig) (gConfig *gorm.Config) {
	dbLogger := &dbLogger{zlog}
	gConfig = &gorm.Config{
		Logger:                                   dbLogger,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	return gConfig
}

// UserRepository ...
type UserRepository interface {
	NewUser(ctx context.Context, user *model.User) (ruser *model.User, err error)
	GetUser(ctx context.Context, uid string) (ruser *model.User, err error)
	GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error)
	UpdateUser(ctx context.Context, user *model.User) (ruser *model.User, err error)
	DeleteUser(ctx context.Context, uid string) (err error)
}

// UserReadOnlyRepository ...
type UserReadOnlyRepository interface {
	GetUser(ctx context.Context, uid string) (ruser *model.User, err error)
	GetUserByEmail(ctx context.Context, email string) (ruser *model.User, err error)
}
