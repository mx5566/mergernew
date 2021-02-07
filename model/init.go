package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	"gorm.io/gorm"
)

var GDB1 *gorm.DB
var GDB2 *gorm.DB
var GDB3 *gorm.DB
var GDB4 *gorm.DB

var (
	DBType        = "mysql"
	DBUser        = "root"
	DBPasswd      = "123456"
	DBHost1       = "127.0.0.1"
	DBHost2       = "127.0.0.1"
	DBNameA       = "game1"
	DBNameB       = "game2"
	DBNameC       = "information_schema"
	DBNameD       = "information_schema"
	DBTablePrefix = ""
)

func NewDB(user, password, host, dbname, tablePrefix string) (*gorm.DB, error) {
	var err error

	//&multiStatements=true
	dia := mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		user, password, host, dbname))
	db, err := gorm.Open(dia, &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "",   // 表名前缀
		SingularTable: true, // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `article`

	},
		Logger:            logger.Default.LogMode(logger.Info),
		AllowGlobalUpdate: true,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
