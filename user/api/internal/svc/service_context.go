package svc

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"pdf-service/user/api/internal/config"
	"pdf-service/user/api/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
	DB        *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := Init(c)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db),
		DB:        db,
	}
}

func Init(c config.Config) (db *gorm.DB) {
	var (
		sqlDB *sql.DB
		err   error
	)
	mysqlConf := mysql.Config{DSN: c.MySQL.DSN}

	gormConfig := configLog(c.MySQL.LogMode)
	if db, err = gorm.Open(mysql.New(mysqlConf), gormConfig); err != nil {
		log.Fatal("opens database failed: ", err)
	}
	if sqlDB, err = db.DB(); err != nil {
		log.Fatal("db.db() failed: ", err)
	}

	sqlDB.SetMaxIdleConns(c.MySQL.MaxIdleCons)
	sqlDB.SetMaxOpenConns(c.MySQL.MaxOpenCons)
	return
}

func configLog(mod bool) (c *gorm.Config) {
	if mod {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 表名不加复数形式，false默认加
			},
		}
	} else {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 表名不加复数形式，false默认加
			},
		}
	}
	return
}
