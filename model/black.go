package model

import "gorm.io/gorm"

// blacklist table struct
type BlackList struct {
	RoleID  uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	BlackID uint32 `gorm:"primaryKey;column:black_id;autoIncrement:false"`
}

func HandleBlackList(db1, db2 *gorm.DB) error {
	var blacks []*BlackList

	err := db2.Table("blacklist").Find(&blacks).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, BlackListC, blacks)
	if err != nil {
		return err
	}

	return nil
}
