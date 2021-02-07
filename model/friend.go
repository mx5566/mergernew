package model

import (
	"gorm.io/gorm"
)

// friend table struct
type Friend struct {
	RoleID   uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	FriendID uint32 `gorm:"primaryKey;column:friend_id;autoIncrement:false"`
	GroupID  int8   `gorm:"column:group_id"`
}

func HandleFriend(db1, db2 *gorm.DB) error {
	var friends []*Friend

	err := db2.Table("friend").Find(&friends).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, FriendC, friends)
	if err != nil {
		return err
	}

	return nil
}
