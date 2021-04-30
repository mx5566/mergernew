package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

// 本机的连接
var GDB1 *gorm.DB
var GDB3 *gorm.DB

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

	//defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(0)
	//sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	return db, nil
}
