package model

import (
	"gorm.io/gorm"
)

// skill table struct
type Skill struct {
	RoleID      uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	ID          uint32 `gorm:"primaryKey;column:id;autoIncrement:false"`
	BiddenLevel int8   `gorm:"not null;column:bidden_level"`
	SelfLevel   int8   `gorm:"not null;column:self_level;default:0"`
	Proficiency int64  `gorm:"not null;column:proficiency;default:0"`
	CoolDown    int32  `gorm:"not null;column:cooldown"`
}

func HandleSkill(db1, db2 *gorm.DB) error {
	var skills []*Skill

	err := db2.Table("skill").Find(&skills).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SkillC, skills)
	if err != nil {
		return err
	}

	return nil
}

// title table struct
type Title struct {
	RoleID uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	ID     uint32 `gorm:"primaryKey;column:id;autoIncrement:false"`
	Count  uint32 `gorm:"not null;column:count"`
	Time   string `gorm:"not null;column:time;char(20)"`
}

func HandleTitle(db1, db2 *gorm.DB) error {
	var ts []*Title

	err := db2.Table("title").Find(&ts).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, TitleC, ts)
	if err != nil {
		return err
	}

	return nil
}
