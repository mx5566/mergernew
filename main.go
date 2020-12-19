package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mx5566/mergernew/model"
	"math"
)

const BaseLength = 10000

func main() {
	model.GDB1, _ = model.NewDB(model.DBType, model.DBUser, model.DBPasswd, model.DBHost, model.DBNameA, model.DBTablePrefix)
	model.GDB2, _ = model.NewDB(model.DBType, model.DBUser, model.DBPasswd, model.DBHost, model.DBNameB, model.DBTablePrefix)
	model.GDB3, _ = model.NewDB(model.DBType, model.DBUser, model.DBPasswd, model.DBHost, model.DBNameC, model.DBTablePrefix)

	model.SetEnv()

	// 远程的数据库数据，合并过来的数据库
	var MinRoleCID uint32
	var MaxRoleCID uint32
	var MaxRoleID uint32
	var increaseNum uint32

	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct


	err := model.GDB1.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MaxRoleID = result.Max


	err = model.GDB2.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MaxRoleCID = result.Max

	err = model.GDB2.Raw("select min(role_id) as min from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MinRoleCID = result.Min

	if MaxRoleID >= MaxRoleCID {
		increaseNum = MaxRoleID + 5
	} else {
		increaseNum = MaxRoleCID + 5
	}
	MinRoleCID += increaseNum

	// 修改所有合并过来表里面有role_id的字段增加increaseNum
	type TablesRoleID struct {
		TableName string `json:"table_name"`
	}

	// 所有有角色ID（role_id）的字段的表
	var tableRoleID []TablesRoleID

	model.GDB3.Table("COLUMNS").Select("table_name as table_name").Where("TABLE_SCHEMA=? and COLUMN_NAME= ?", model.DBNameB, "role_id").Find(&tableRoleID)

	fmt.Println(tableRoleID)

	// 更新所有上面有role_id的表
	/*for _, table := range tableRoleID {
		model.GDB2.Table(table.TableName).Update("role_id", gorm.Expr("role_id + ?", increaseNum))
	}*/

	//把B表里面所有角色ID需要修改的也增加increaseNum
	err = model.GDB2.Table("item").AddIndex("temp_item_o_c", "create_mode").Error
	if err != nil {
		panic(err)
	}

	type ItemIDs struct {
		ID int64 `gorm:"column:serial"`
	}

	var itemIDs []ItemIDs
	err = model.GDB2.Table("item").Select("serial as serial").Find(&itemIDs).Error
	if err != nil {
		panic(err)
	}

	//model.GDB2.Begin() 好像加了事务还慢了
	length := len(itemIDs)

	l1 := math.Ceil(float64(length * 1.0 / (1.0 * BaseLength)))

	fmt.Println(l1)
	model.GDB2.Begin()
	err = model.GDB2.Table("item").Update(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		model.GDB2.Rollback()
		panic(err)
	}
	model.GDB2.Commit()

	// 所有mode为8的更新increaseNum
	model.GDB2.Begin()
	err = model.GDB2.Table("item").Where("create_mode = ?", 8).Update(map[string]interface{}{"create_id": gorm.Expr("create_id + ?", increaseNum)}).Error
	if err != nil {
		model.GDB2.Rollback()
		panic(err)
	}
	model.GDB2.Commit()

	//移除索引
	err = model.GDB2.Table("item").RemoveIndex("temp_item_o_c").Error
	if err != nil {
		panic(err)
	}






	defer model.GDB1.Close()
	defer model.GDB2.Close()
}