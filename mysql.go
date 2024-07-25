package gormHelper

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Mysql struct {
	Databases []MysqlDatabase
}

type MysqlDatabase struct {
	Name            string
	Dsn             string
	Env             string
	MysqlConfig     *mysql.Config
	GormConfig      *gorm.Config
	MaxIdleConns    *int
	MaxOpenConns    *int
	ConnMaxLifetime *time.Duration
}

func (m *Mysql) Get() map[string]MysqlDatabase {
	databaseMap := make(map[string]MysqlDatabase, len(m.Databases))
	for _, database := range m.Databases {
		//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
		databaseMap[database.Name] = database
	}
	return databaseMap
}

func (database *MysqlDatabase) Db() (*gorm.DB, error) {
	mysqlConfig := *database.MysqlConfig
	if database.MysqlConfig != nil {
		mysqlConfig.DSN = database.Dsn
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), database.GormConfig); err != nil {
		return db, err
	} else {
		if database.Env == "debug" {
			db.Debug()
		}
		var sqlDb *sql.DB
		if sqlDb, err = db.DB(); err != nil {
			return db, err
		} else {
			if database.MaxIdleConns != nil {
				sqlDb.SetMaxIdleConns(*database.MaxIdleConns)
			}
			if database.MaxOpenConns != nil {
				sqlDb.SetMaxOpenConns(*database.MaxOpenConns)
			}
			if database.ConnMaxLifetime != nil {
				sqlDb.SetConnMaxLifetime(*database.ConnMaxLifetime)
			}
		}
		return db, nil
	}
}

func (m *Mysql) Paginate(db *gorm.DB, page int, size int) *gorm.DB {
	if page == 0 {
		page = 1
	}
	if size > 100 {
		size = 100
	}
	if size <= 0 {
		size = 1
	}
	return db.Offset((page - 1) * size).Limit(size)
}