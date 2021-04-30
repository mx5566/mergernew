package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
)

type JimaiItem struct {
	ID      int64 `gorm:"primary_key;column:id"`
	PayType int16 `gorm:"primary_key;column:pay_type"`
	Money   int64 `gorm:"primary_key;column:money"`
	Time    int64 `gorm:"primary_key;column:time"`
}

func HandleJiMai(db1, db2 *gorm.DB, mapItems map[int64]*ItemEx) error {
	var jiMaiItems []*JimaiItem
	err := db2.Table("jimai_item").Select("id, pay_type, money, time").Find(&jiMaiItems).Error
	if err != nil {
		return err
	}

	deleteIndex := make([]int, 0)
	// 处理装备的ID
	for key, value := range jiMaiItems {
		if v, ok := mapItems[value.ID]; ok {
			// 把新的id给装备
			jiMaiItems[key].ID = v.Serial
		} else {
			//
			deleteIndex = append(deleteIndex, key)
			logm.ErrorfE("jiMaiItems[%d] not found in item table", value.ID)
		}
	}

	var count = 0
	for _, v := range deleteIndex {
		jiMaiItems = append(jiMaiItems[0:v-count], jiMaiItems[v+1-count:]...)
		count++
	}

	// 批量插入被合数据库表item
	err = BatchSave(db1, JiMaiItemC, jiMaiItems)

	if err != nil {
		return err
	}

	jiMaiItems = make([]*JimaiItem, 0)

	return nil
}
