package mysql

import (
	"exchange/services/asset-service/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-micro.dev/v4/logger"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

func New(conf *config.Config) *Repository {
	repo := &Repository{DB: newDB(conf)}
	return repo
}
func newDB(conf *config.Config) *gorm.DB {
	db, err := gorm.Open("mysql", conf.SqlMap["community"].DSN)
	if err != nil || db == nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(conf.SqlMap["community"].MaxIdle)
	db.DB().SetMaxOpenConns(conf.SqlMap["community"].MaxConn)
	db.DB().SetConnMaxLifetime(time.Duration(conf.SqlMap["community"].MaxLifeTime))
	db.LogMode(conf.Log.MysqlLog)
	if err := db.DB().Ping(); err != nil {
		panic(err)
	} else {
		logger.Infof("community db connected.")
	}
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	return db
}

// updateTimeStampForCreateCallback
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now().UnixNano() / 1000 / 1000
		if createdAtField, ok := scope.FieldByName("create_time"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("update_time"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}
func (r *Repository) Close() {
	_ = r.DB.Close()
}