package model

import "gorm.io/gorm"

// role_extend table struct
type RoleExtend struct {
	RoleID           uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	RoleExtendValue  []byte `gorm:"column:role_extend_value"`
	RoleEquipPosInfo []byte `gorm:"column:role_equip_pos_info"`
}

func HandleRoleExtend(db1, db2 *gorm.DB) error {
	var roleExtends []*RoleExtend

	err := db2.Table("enemy").Find(&roleExtends).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, RoleExtendC, roleExtends)
	if err != nil {
		return err
	}

	return nil
}
