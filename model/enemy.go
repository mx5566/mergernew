package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
)

// enemy table struct
type Enemy struct {
	RoleID   uint32 `gorm:"primaryKey;column:role_id;autoIncrement:false"`
	EnemyID  uint32 `gorm:"primaryKey;column:enemy_id;autoIncrement:false"`
	KillTime int64  `gorm:"not null;column:kill_time;default:0"`
}

func HandleEnemy(db1, db2 *gorm.DB, mapRoleIDsRf map[uint32]uint32) error {
	var enemys1 []*Enemy
	var enemys2 []*Enemy

	err := db1.Table("enemy").Find(&enemys1).Error
	if err != nil {
		return err
	}

	err = db2.Table("enemy").Find(&enemys2).Error
	if err != nil {
		return err
	}

	/////////////////////////////////////////////////////
	deleteIndex := make([]int, 0)
	for i, value := range enemys1 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 敌人表角色本身被隐藏了 清理敌人表这一行数据 roleID[%d] targetID[%d]", value.RoleID, value.EnemyID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.EnemyID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 敌人表敌人本身被隐藏了 清理敌人表这一行数据 roleID[%d] targetID[%d]", value.RoleID, value.EnemyID)
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		enemys1 = append(enemys1[0:value-count], enemys1[value+1-count:]...)
		count++
	}

	deleteIndex = make([]int, 0)
	for i, value := range enemys2 {
		if _, ok := mapRoleIDsRf[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 敌人表角色本身被隐藏了 清理敌人表这一行数据 roleID[%d] targetID[%d]", value.RoleID, value.EnemyID)
			continue
		}

		if _, ok := mapRoleIDsRf[value.EnemyID]; ok {
			deleteIndex = append(deleteIndex, i)

			logm.WarnfE("合区 敌人表敌人本身被隐藏了 清理敌人表这一行数据 roleID[%d] targetID[%d]", value.RoleID, value.EnemyID)
		}
	}

	count = 0
	for _, value := range deleteIndex {
		enemys2 = append(enemys2[0:value-count], enemys2[value+1-count:]...)
		count++
	}

	enemys1 = append(enemys1, enemys2...)

	enemys2 = make([]*Enemy, 0)

	err = BatchEnemySave(db1, enemys1)
	if err != nil {
		return err
	}

	/////////////////////////////////////////////////////

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	//err = BatchSave(db1, EnemyC, enemys)
	//if err != nil {
	//	return err
	//}

	return nil
}

func BatchEnemySave(db1 *gorm.DB, enemys []*Enemy) error {
	err := db1.Exec("truncate table enemy;").Error
	if err != nil {
		return err
	}

	err = BatchSave(db1, EnemyC, enemys)
	if err != nil {
		return err
	}

	return nil
}
