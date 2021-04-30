package model

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/config"
	"gorm.io/gorm"
	"time"
)

// 各种其他的表 直接sql语句处理
func HandleRelation(db1, db2, db3, db4 *gorm.DB, increaseNum uint32) error {
	// 修改所有合并过来表里面有role_id的字段增加increaseNum
	type TablesRoleID struct {
		TableName string `json:"table_name"`
	}

	// 所有有角色ID（role_id）的字段的表
	var tableRoleID []TablesRoleID

	ms := config.Config.Mysql

	db3.Table("COLUMNS").Select("table_name as table_name").Where("TABLE_SCHEMA=? and COLUMN_NAME= ?", ms.DBNameB, "role_id").Find(&tableRoleID)

	t1 := time.Now().Unix()
	// 更新所有上面有role_id的表
	for _, table := range tableRoleID {
		err := db2.Table(table.TableName).Order("role_id desc").Update("role_id", gorm.Expr("role_id + ?", increaseNum)).Error
		if err != nil {
			return err
		}
	}

	logm.DebugfE("所有有角色id的表都处理了 耗时[%d]", time.Now().Unix()-t1)
	t1 = time.Now().Unix()

	/*
		UPDATE guild_hebin_copy SET creater_name_id=creater_name_id+increase_num;
		UPDATE guild_hebin_copy SET leader_id=leader_id+increase_num;
		UPDATE item_del_hebin_copy SET owner_id=owner_id+increase_num;
	*/
	err := db2.Table("guild").Updates(map[string]interface{}{"creater_name_id": gorm.Expr("creater_name_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	err = db2.Table("guild").Updates(map[string]interface{}{"leader_id": gorm.Expr("leader_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	logm.DebugfE("处理了工会表里面与角色id有关的 耗时[%d]", time.Now().Unix()-t1)
	t1 = time.Now().Unix()

	err = db2.Table("item_del").Updates(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	/*
		update friend_hebin_copy SET friend_id = friend_id + increase_num;
		update blacklist_hebin_copy set black_id = black_id + increase_num;
		update enemy_hebin_copy set enemy_id = enemy_id + increase_num;
	*/

	err = db2.Table("friend").Order("friend_id desc").Updates(map[string]interface{}{"friend_id": gorm.Expr("friend_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	logm.DebugfE("处理了好友表里面与角色id有关的 耗时[%d]", time.Now().Unix()-t1)
	t1 = time.Now().Unix()

	err = db2.Table("blacklist").Order("black_id desc").Updates(map[string]interface{}{"black_id": gorm.Expr("black_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	logm.DebugfE("处理了黑名单表里面与角色id有关的 耗时[%d]", time.Now().Unix()-t1)
	t1 = time.Now().Unix()

	err = db2.Table("enemy").Order("enemy_id desc").Updates(map[string]interface{}{"enemy_id": gorm.Expr("enemy_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	logm.DebugfE("处理了敌人表里面与角色id有关的 耗时[%d]", time.Now().Unix()-t1)
	t1 = time.Now().Unix()

	/*
		update mail_hebin_copy set recv_role_id = recv_role_id + increase_num where recv_role_id != 4294967295;
		update mail_hebin_copy set send_role_id = send_role_id + increase_num where send_role_id != 4294967295;
	*/

	err = db2.Table("mail").Where("recv_role_id != ?", 4294967295).Updates(map[string]interface{}{"recv_role_id": gorm.Expr("recv_role_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	err = db2.Table("mail").Where("send_role_id != ?", 4294967295).Updates(map[string]interface{}{"send_role_id": gorm.Expr("send_role_id + ?", increaseNum)}).Error
	if err != nil {
		return err
	}

	logm.DebugfE("处理了邮件表里面与角色id有关的 耗时[%d]", time.Now().Unix()-t1)

	return nil
}
