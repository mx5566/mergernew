package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var GDB1 *gorm.DB
var GDB2 *gorm.DB
var GDB3 *gorm.DB

var (
DBType = "mysql"
DBUser = "root"
DBPasswd = "123456"
DBHost = "127.0.0.1"
DBNameA = "game1"
DBNameB = "game2"
DBNameC = "information_schema"
DBTablePrefix = ""
)

func NewDB(dbType, user, password, host, dbname, tablePrefix string) (*gorm.DB, error){
	var err error

	db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, dbname))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return tablePrefix + defaultTableName;
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	return db, nil
}

func TestItem() {
	//把B的数据插入到A里面
	var items []Item
	var itemFields = []string{"num", "type_id", "bind", "lock_state", "use_times", "first_gain_time", "create_mode",
		"create_id", "creator_id", "create_time", "owner_id", "account_id", "container_type_id", "suffix", "name_id", "bind_time",
		"script_data1", "script_data2", "create_bind", "strdwExternData", "serial", "'_copy'"}
	err := GDB2.Select(itemFields).Find(&items).Error

	if err != nil {
		fmt.Println(err.Error())
	}
}
