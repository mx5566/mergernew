package model

import "gorm.io/gorm"

// enemy table struct
type Enemy struct {
	RoleID  uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	EnemyID uint32 `gorm:"primaryKey;column:enemy_id;autoIncrement:false"`
	BlackID uint32 `gorm:"not null;column:kill_time;default:0"`
}

func HandleEnemy(db1, db2 *gorm.DB) error {
	var enemys []*Enemy

	err := db2.Table("enemy").Find(&enemys).Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, EnemyC, enemys)
	if err != nil {
		return err
	}

	return nil
}
