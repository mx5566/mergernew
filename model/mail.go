package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
)

type Mail struct {
	MailID        uint32 `gorm:"primary_key;column:mail_id"`
	MailName      []byte `gorm:"column:mail_name"`
	MailContent   []byte `gorm:"column:mail_content"`
	SendRoleID    uint32 `gorm:"column:send_role_id;default:4294967295"`
	RecvRoleID    uint32 `gorm:"column:recv_role_id;default:4294967295"`
	IsSend        uint8  `gorm:"column:is_send"`
	IsRead        uint8  `gorm:"column:is_read"`
	IsWidthDrawal uint8  `gorm:"column:is_withdrawal"`
	Solve         uint32 `gorm:"column:solve"`
	GiveMoney     uint32 `gorm:"column:give_money"`
	ItemSerial    []byte `gorm:"column:item_serial"`
	SendTime      string `gorm:"column:send_time;char(20)"`
	RecvTime      string `gorm:"column:recv_time;char(20)"`
	IsAtOnce      uint8  `gorm:"column:is_at_once"`
	MoneyType     uint8  `gorm:"column:moneytype"`
	YuanBaoType   uint8  `gorm:"column:yuanbao_type"`
	ItemSerialNew int64  `gorm:"column:item_serial_new"`
}

// 这个可以在物品的所有逻辑处理之前就处理掉
// 后面要把物品的数据加载进内存里面处理
// 加载进内存的处理方式如果数据量很大的话，也是一个问题，内存不够用了。。。哎哎哎
func HandleMail(db1, db2 *gorm.DB) (error, uint32) {
	// 远程的数据库数据，合并过来的数据库
	var MAXMailID uint32
	var increaseNum uint32

	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct

	err := db1.Raw("select max(mail_id) as max from mail;").Scan(&result).Error
	if err != nil {
		return err, 0
	}
	MAXMailID = result.Max
	increaseNum = MAXMailID + 5

	// update mail_hebin_copy set mail_id = mail_id + increase_num;
	err = db2.Table("mail").Update("mail_id", gorm.Expr("mail_id + ?", increaseNum)).Error
	if err != nil {
		return err, 0
	}

	return nil, increaseNum
}

func HandleMailByIncreaseID(db1, db2 *gorm.DB) (map[uint32]uint32, error) {
	var MAXMailID uint32
	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct

	err := db1.Raw("select max(mail_id) as max from mail;").Scan(&result).Error
	if err != nil {
		return nil, err
	}
	MAXMailID = result.Max
	increaseNum := MAXMailID

	// id采用自增长id的方式去处理，避免跳跃的增长导致32为的不够用
	var mails []*MailEx
	err = db2.Table("mail").Select("mail_id").Order("mail_id asc").Find(&mails).Error
	if err != nil {
		return nil, err
	}

	mapIDs := make(map[uint32]uint32)
	for _, v1 := range mails {
		increaseNum++
		mapIDs[v1.MailID] = increaseNum
	}

	return mapIDs, nil
}

func HandleMailItem(db1, db2 *gorm.DB, mapItems map[int64]*ItemEx, mapMails map[uint32]uint32) error {
	var mailItems []*Mail
	err := db2.Table("mail").Find(&mailItems).Error
	if err != nil {
		return err
	}

	deleteIndex := make([]int, 0)
	// 处理物品的ID
	for key, value := range mailItems {
		need := false
		if v, ok := mapMails[value.MailID]; ok {
			mailItems[key].MailID = v
		} else {
			//deleteIndex = append(deleteIndex, key)
			need = true
		}

		if value.ItemSerialNew == 0 || value.ItemSerialNew == 4294967295 {
			continue
		}

		if v, ok := mapItems[value.ItemSerialNew]; ok {
			// 把新的id给mail
			mailItems[key].ItemSerialNew = v.Serial
		} else {
			//
			need = true
			logm.ErrorfE("MailItems[%d] not found in item table", value.ItemSerialNew)
		}

		if need {
			deleteIndex = append(deleteIndex, key)
		}
	}

	var count = 0
	for _, v := range deleteIndex {
		mailItems = append(mailItems[0:v-count], mailItems[v+1-count:]...)
		count++
	}

	// 批量插入被合数据库表item
	err = BatchSave(db1, MailItemC, mailItems)
	if err != nil {
		return err
	}

	mailItems = make([]*Mail, 0)

	return nil
}
