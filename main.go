package main

import (
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/model"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {

	t1 := time.Now().Unix()

	model.GDB1, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameA, model.DBTablePrefix)
	model.GDB2, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameB, model.DBTablePrefix)
	model.GDB3, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameC, model.DBTablePrefix)

	model.SetEnv()
	logm.Init("mergernew", map[string]string{"errFile": "game_err.log", "logFile": "game.log"}, "debug")

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

	//fmt.Println(tableRoleID)

	// 更新所有上面有role_id的表
	for _, table := range tableRoleID {
		model.GDB2.Table(table.TableName).Update("role_id", gorm.Expr("role_id + ?", increaseNum))
	}

	// 用来测试时间
	// 不利用主键去count 利用第二索引去count
	var c1 int64 = 0
	err = model.GDB2.Table("item").Select("count(container_type_id) as count").Where("container_type_id > ?", 0).Count(&c1).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("最大个数 ", c1)

	//把B表里面所有角色ID需要修改的也增加increaseNum

	/*type ItemIDs struct {
		ID int64 `gorm:"column:serial"`
	}

	var itemIDs []ItemIDs
	err = model.GDB2.Table("item").Select("serial as serial").Find(&itemIDs).Error
	if err != nil {
		panic(err)
	}*/

	/////////////////////
	// owner_id	owner_id	NORMAL	0	A	40			0
	//model.GDB2.Begin() 好像加了事务还慢了
	//model.GDB2.Begin()
	//移除索引
	/*err = model.GDB2.Exec("ALTER TABLE item DROP INDEX owner_id ;").Error
	if err != nil {
		panic(err)
	}*/

	/*err = model.GDB2.Table("item").Where("container_type_id in ?", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15, 16}).Updates(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		//model.GDB2.Rollback()
		panic(err)
	}*/

	// 删除索引
	// 删除物品表里面的两个索引 提升update的性能
	//err = model.GDB2.Exec("alter table item drop index account_id; alter table item drop index owner_id;").Error
	//err = model.GDB2.Exec("alter table item drop index account_id; alter table item drop index container_type_id;alter table item drop index owner_id;").Error
	//if err != nil {
	//	panic(err)
	//}

	// handle mail 先处理掉
	err = model.HandleMail(model.GDB1, model.GDB2)
	if err != nil {
		panic(err)
	}

	// handle item relative
	err = model.HandleItemRelation(model.GDB1, model.GDB2, increaseNum)
	if err != nil {
		panic(err)
	}

	//err = model.HandleItemOwnerId(model.GDB1, model.GDB2, increaseNum)
	//if err != nil {
	//	panic(err)
	//}

	// 增加删除的索引
	//err = model.GDB2.Exec("alter table item add index account_id (account_id, container_type_id);" +
	/*"alter table item add index container_type_id (`container_type_id`);"*/
	//	"alter table item add index owner_id (`owner_id`);").Error

	//model.GDB2.Commit()

	/////////////////////

	/////////////////////

	// 所有mode为8的更新increaseNum
	//model.GDB2.Begin()
	err = model.GDB2.Exec("ALTER TABLE `item` ADD INDEX temp_item_o_c (`create_mode`) ").Error
	if err != nil {
		//model.GDB2.Rollback()
		panic(err)
	}

	err = model.GDB2.Table("item").Where("create_mode = ?", 8).Updates(map[string]interface{}{"create_id": gorm.Expr("create_id + ?", increaseNum)}).Error
	if err != nil {
		//model.GDB2.Rollback()
		panic(err)
	}
	//model.GDB2.Commit()

	//移除索引
	err = model.GDB2.Exec("alter table item drop index temp_item_o_c ;").Error
	if err != nil {
		panic(err)
	}
	/////////////////////

	/*
		UPDATE guild_hebin_copy SET creater_name_id=creater_name_id+increase_num;
		UPDATE guild_hebin_copy SET leader_id=leader_id+increase_num;
		UPDATE item_del_hebin_copy SET owner_id=owner_id+increase_num;
	*/
	err = model.GDB2.Table("guild").Updates(map[string]interface{}{"creater_name_id": gorm.Expr("creater_name_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("guild").Updates(map[string]interface{}{"leader_id": gorm.Expr("leader_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("item_del").Updates(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	/*
		update friend_hebin_copy SET friend_id = friend_id + increase_num;
		update blacklist_hebin_copy set black_id = black_id + increase_num;
		update enemy_hebin_copy set enemy_id = enemy_id + increase_num;
	*/

	err = model.GDB2.Table("friend").Updates(map[string]interface{}{"friend_id": gorm.Expr("friend_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("blacklist").Updates(map[string]interface{}{"black_id": gorm.Expr("black_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("enemy").Updates(map[string]interface{}{"enemy_id": gorm.Expr("enemy_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	/*
		update mail_hebin_copy set recv_role_id = recv_role_id + increase_num where recv_role_id != 4294967295;
		update mail_hebin_copy set send_role_id = send_role_id + increase_num where send_role_id != 4294967295;
	*/

	err = model.GDB2.Table("mail").Where("recv_role_id != ?", 4294967295).Updates(map[string]interface{}{"recv_role_id": gorm.Expr("recv_role_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("mail").Where("send_role_id != ?", 4294967295).Updates(map[string]interface{}{"send_role_id": gorm.Expr("send_role_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.HandleRoleNameOrigin(model.GDB1, model.GDB2, model.GDB3)
	if err != nil {
		panic(err)
	}

	// handle account_common
	err = model.HandleAccountCommon(model.GDB1, model.GDB2)
	if err != nil {
		panic(err)
	}

	// add 2 item to 1
	// 把合并的物品数据插入到被合并数据库物品表
	//err = model.HandleItem(model.GDB1, model.GDB2)
	//if err != nil {
	//	fmt.Println(err)
	//panic(err)
	//}

	fmt.Println("消耗的时间", time.Now().Unix()-t1)
}
