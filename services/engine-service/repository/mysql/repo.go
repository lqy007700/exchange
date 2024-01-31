package mysql

import (
	"engine-service/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-micro.dev/v4/logger"
	"time"
)

func newDB() *gorm.DB {
	conf := config.Conf
	logger.Info(conf.SqlMap)
	db, err := gorm.Open("mysql", conf.SqlMap["asset"].DSN)
	if err != nil || db == nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(conf.SqlMap["asset"].MaxIdle)
	db.DB().SetMaxOpenConns(conf.SqlMap["asset"].MaxConn)
	db.DB().SetConnMaxLifetime(time.Duration(conf.SqlMap["asset"].MaxLifeTime))
	db.LogMode(conf.Log.MysqlLog)
	if err := db.DB().Ping(); err != nil {
		panic(err)
	} else {
		logger.Infof("asset db connected.")
	}
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
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
