package model

import (
	"gorm.io/gorm"
)

// map_limit table struct
type MapLimit struct {
	ID       uint32 `gorm:"primaryKey;column:id;autoIncrement"`
	RoleID   uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	MapID    uint32 `gorm:"primaryKey;column:map_id;autoIncrement:false"`
	Type     uint8  `gorm:"column:type;not null"`
	EnterNum uint32 `gorm:"column:enter_num;not null"`
}

func HandleMapLimit(db1, db2 *gorm.DB) error {
	var limits []*MapLimit

	err := db2.Table("map_limit").Find(&limits).Error
	if err != nil {
		return err
	}

	for k1, _ := range limits {
		limits[k1].ID = 0
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, MapLimitC, limits)
	if err != nil {
		return err
	}

	return nil
}
