package model

import "gorm.io/gorm"

// buff table struct
type Buff struct {
	RoleID        uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	SrcUnitID     uint32 `gorm:"column:src_unit_id"`
	SrcSkillID    uint32 `gorm:"column:src_skill_id"`
	ItemTypeID    uint32 `gorm:"column:item_type_id"`
	ItemSerialID  int64  `gorm:"column:item_serial_id"`
	BuffID        uint32 `gorm:"primaryKey;column:buff_id;autoIncrement:false"`
	CurTick       uint32 `gorm:"column:cur_tick"`
	Level         int8   `gorm:"column:level"`
	CurLapTimes   int8   `gorm:"column:cur_lap_times"`
	EffectSkillID []byte `gorm:"column:effect_skill_id"`
}

func HandleBuff(db1, db2 *gorm.DB) error {
	var buffs []*Buff

	err := db2.Table("buff").Find(&buffs).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, BuffC, buffs)
	if err != nil {
		return err
	}

	return nil
}
