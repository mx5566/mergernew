package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
)

// blacklist table struct
type BlackList struct {
	RoleID  uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	BlackID uint32 `gorm:"primaryKey;column:black_id;autoIncrement:false"`
}

func HandleBlackList(db1, db2 *gorm.DB, mapRoleIDsRf map[uint32]uint32) error {
	var blacks1 []*BlackList
	var blacks2 []*BlackList

	err := db1.Table("blacklist").Find(&blacks1).Error
	if err != nil {
		return err
	}

	err = db2.Table("blacklist").Find(&blacks2).Error
	if err != nil {
		return err
	}

	////////////////////////////////////////////
	// 处理合区被隐藏的关系做清理
	deleteIndex := make([]int, 0)
	for i, value := range blacks1 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 黑名单角色被隐藏 清理关系 roleID[%d] targetID[%d]", value.RoleID, value.BlackID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.BlackID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 黑名单目标被隐藏 清理关系 roleID[%d] targetID[%d]", value.RoleID, value.BlackID)
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		blacks1 = append(blacks1[0:value-count], blacks1[value+1-count:]...)
		count++
	}

	deleteIndex = make([]int, 0)
	for i, value := range blacks2 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 黑名单角色被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.BlackID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.BlackID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 黑名单目标被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.BlackID)
		}
	}

	count = 0
	for _, value := range deleteIndex {
		blacks2 = append(blacks2[0:value-count], blacks2[value+1-count:]...)
		count++
	}

	blacks1 = append(blacks1, blacks2...)

	blacks2 = make([]*BlackList, 0)
	////////////////////////////////////////////

	err = BatchBlackListSave(db1, blacks1)
	if err != nil {
		return err
	}
	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	//err = BatchSave(db1, BlackListC, blacks2)
	//if err != nil {
	//	return err
	//}

	return nil
}

func BatchBlackListSave(db1 *gorm.DB, blacks []*BlackList) error {
	err := db1.Exec("truncate table blacklist;").Error
	if err != nil {
		return err
	}

	err = BatchSave(db1, BlackListC, blacks)
	if err != nil {
		return err
	}

	return nil
}
