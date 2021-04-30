package model

import "gorm.io/gorm"

// reward table struct
type Reward struct {
	RoleID     uint32 `gorm:"primaryKey;column:role_id;autoIncrement"`
	RewardData []byte `gorm:"column:reward_data"`
}

func HandleReward(db1, db2 *gorm.DB) error {
	var rewards []*Reward

	err := db2.Table("reward").Find(&rewards).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, RewardC, rewards)
	if err != nil {
		return err
	}

	return nil
}

// sign table struct
type Sign struct {
	RoleID       uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	Data         []byte `gorm:"column:data"`
	RewardData   []byte `gorm:"column:reward_data"`
	MianQianTime uint32 `gorm:"not null;column:mianqian_time"`
}

func HandleSign(db1, db2 *gorm.DB) error {
	var signs []*Sign

	err := db2.Table("sign").Find(&signs).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SignC, signs)
	if err != nil {
		return err
	}

	return nil
}
