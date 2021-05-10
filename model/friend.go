package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
)

// friend table struct
type Friend struct {
	RoleID   uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	FriendID uint32 `gorm:"primaryKey;column:friend_id;autoIncrement:false"`
	GroupID  int8   `gorm:"column:group_id"`
}

func HandleFriend(db1, db2 *gorm.DB, mapRoleIDsRf map[uint32]uint32) error {
	var friends1 []*Friend
	var friends2 []*Friend

	err := db1.Table("friend").Find(&friends1).Error
	if err != nil {
		return err
	}

	err = db2.Table("friend").Find(&friends2).Error
	if err != nil {
		return err
	}

	///////////////////////////////////////////////
	deleteIndex := make([]int, 0)
	for i, value := range friends1 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 好友表自己角色被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.FriendID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.FriendID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 好友表好友角色被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.FriendID)
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		friends1 = append(friends1[0:value-count], friends1[value+1-count:]...)
		count++
	}

	deleteIndex = make([]int, 0)
	for i, value := range friends2 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 好友表自己角色被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.FriendID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.FriendID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 好友表好友角色被隐藏 清理好友关系 roleID[%d] targetID[%d]", value.RoleID, value.FriendID)
		}
	}

	count = 0
	for _, value := range deleteIndex {
		friends2 = append(friends2[0:value-count], friends2[value+1-count:]...)
		count++
	}

	friends1 = append(friends1, friends2...)

	friends2 = make([]*Friend, 0)

	err = BatchFriendSave(db1, friends1)
	if err != nil {
		return err
	}
	///////////////////////////////////////////////

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	//err = BatchSave(db1, FriendC, friends)
	//if err != nil {
	//	return err
	//}

	return nil
}

func BatchFriendSave(db1 *gorm.DB, friends []*Friend) error {
	err := db1.Exec("truncate table friend;").Error
	if err != nil {
		return err
	}

	err = BatchSave(db1, FriendC, friends)
	if err != nil {
		return err
	}

	return nil
}
