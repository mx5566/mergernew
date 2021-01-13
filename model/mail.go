package model

import "gorm.io/gorm"

// 这个可以在物品的所有逻辑处理之前就处理掉
// 后面要把物品的数据加载进内存里面处理
// 加载进内存的处理方式如果数据量很大的话，也是一个问题，内存不够用了。。。哎哎哎
func HandleMail(db1, db2 *gorm.DB) error {
	// 远程的数据库数据，合并过来的数据库
	var MAXMailID uint32
	var MAXMailCID uint32
	var increaseNum uint32

	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct

	err := db1.Raw("select max(mail_id) as max from mail;").Scan(&result).Error
	if err != nil {
		return err
	}
	MAXMailID = result.Max

	err = db2.Raw("select max(mail_id) as max from mail;").Scan(&result).Error
	if err != nil {
		return err
	}
	MAXMailCID = result.Max

	if MAXMailID >= MAXMailCID {
		increaseNum = MAXMailID + 5
	} else {
		increaseNum = MAXMailCID + 5
	}

	// update mail_hebin_copy set mail_id = mail_id + increase_num;
	err = db2.Table("mail").Update("mail_id", gorm.Expr("mail_id + ?", increaseNum)).Error
	if err != nil {
		return err
	}

	// update item_hebin_copy set owner_id = owner_id + increase_num where container_type_id = 11;
	err = db2.Table("item").Where("container_type_id = ?", 11).Update("owner_id", gorm.Expr("owner_id + ?", increaseNum)).Error
	if err != nil {
		return err
	}

	return nil
}
