package model

import "gorm.io/gorm"

// left_msg table struct
type LeftMsg struct {
	MsgID   uint32 `gorm:"primaryKey;column:msg_id;autoIncrement"`
	RoleID  uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	MsgData []byte `gorm:"column:msg_data"`
}

func HandleLeftMsg(db1, db2 *gorm.DB) error {
	var msgs []*LeftMsg

	err := db2.Table("left_msg").Find(&msgs).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, LeftMsgC, msgs)
	if err != nil {
		return err
	}

	return nil
}
