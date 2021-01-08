package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"math"
	"sync"
)

const BaseLength = 1000
const BaseLength1000 = 100

const (
	Ac = iota + 1
	ItemC
	ItemOwnerID
)

// BatchSave 批量插入数据
func BatchSave(db *gorm.DB, t int, arr interface{}) error {
	switch t {
	case Ac:
		accountCommon := arr.([]*AccountCommon)

		var buffer bytes.Buffer
		length := len(accountCommon)
		count := math.Ceil(float64(length) / float64(BaseLength))
		for i := 0; i < int(count); i++ {
			buffer.Reset()
			end := int(math.Min(float64((i+1)*BaseLength), float64(length)))
			start := i * BaseLength

			sql := "INSERT INTO `account_common`(`account_id`, `account_name`, `safecode_crc`, `reset_time`, `bag_password_crc`, `baibao_yuanbao`, `ware_size`, `ware_silver`, `warestep`, `yuanbao_recharge`, `IsReceive`, `total_recharge`, `receive_type`, `receive_type_ex`, `web_type`, `score`, `robort`, `forbid_flag`, `data_ex`) VALUES"
			if _, err := buffer.WriteString(sql); err != nil {
				return err
			}

			for j := start; j < end; j++ {
				if j == end-1 {
					buffer.WriteString(fmt.Sprintf("(%d,'%s',%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, '%s');", accountCommon[j].AccountID, accountCommon[j].AccountName, accountCommon[j].SafeCodeCrc, accountCommon[j].ResetTime, accountCommon[j].BagPasswordCrc, accountCommon[j].BaiBaoYuanBao, accountCommon[j].WareSize, accountCommon[j].WareSilver, accountCommon[j].WareStep, accountCommon[j].YuanBaoRecharge, accountCommon[j].IsReceive, accountCommon[j].TotalRecharge, accountCommon[j].ReceiveType, accountCommon[j].ReceiveTypeEx, accountCommon[j].WebType, accountCommon[j].Score, accountCommon[j].Robort, accountCommon[j].Forbid_flag, accountCommon[j].DataEx))
				} else {
					buffer.WriteString(fmt.Sprintf("(%d,'%s',%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, '%s'),", accountCommon[j].AccountID, accountCommon[j].AccountName, accountCommon[j].SafeCodeCrc, accountCommon[j].ResetTime, accountCommon[j].BagPasswordCrc, accountCommon[j].BaiBaoYuanBao, accountCommon[j].WareSize, accountCommon[j].WareSilver, accountCommon[j].WareStep, accountCommon[j].YuanBaoRecharge, accountCommon[j].IsReceive, accountCommon[j].TotalRecharge, accountCommon[j].ReceiveType, accountCommon[j].ReceiveTypeEx, accountCommon[j].WebType, accountCommon[j].Score, accountCommon[j].Robort, accountCommon[j].Forbid_flag, accountCommon[j].DataEx))
				}
			}

			str := buffer.String()
			fmt.Println("批量保存 字节长度account_common: ", len(str), "byte")
			// 批量更新account_common
			err := db.Exec(str).Error
			if err != nil {
				return err
			}
		}
	case ItemC:
		items := arr.([]*Item)

		err := db.Session(&gorm.Session{PrepareStmt: false, CreateBatchSize: 100}).Create(items).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func BatchUpdate(db *gorm.DB, t int, arr interface{}, args ...map[string]interface{}) error {
	switch t {
	case Ac:
		alreadyAccount := arr.([]*AccountCommon)

		var buffer bytes.Buffer

		repeatLength := len(alreadyAccount)
		count := math.Ceil(float64(repeatLength) / float64(BaseLength))

		// 对于存在的账号数据叠加
		for i := 0; i < int(count); i++ {
			buffer.Reset()

			end := int(math.Min(float64((i+1)*BaseLength), float64(repeatLength)))
			start := i * BaseLength
			for j := start; j < end; j++ {
				buffer.WriteString(fmt.Sprintf("update account_common set baibao_yuanbao=%d, total_recharge=%d, yuanbao_recharge=%d, data_ex='%s' where account_id=%d;",
					alreadyAccount[j].BaiBaoYuanBao, alreadyAccount[j].TotalRecharge, alreadyAccount[j].YuanBaoRecharge, alreadyAccount[j].DataEx, alreadyAccount[j].AccountID))
			}

			str := buffer.String()
			fmt.Println("批量更新 字节长度 account_common: ", len(str), "byte")
			// 批量更新account_common
			err := db.Exec(str).Error
			if err != nil {
				return err
			}
		}
	case ItemOwnerID:
		items := arr.([]*Item1)

		c, _ := args[0]["count"]

		//var buffer bytes.Buffer

		repeatLength := len(items)
		count := math.Ceil(float64(repeatLength) / float64(BaseLength))

		var wg sync.WaitGroup
		wg.Add(int(count))
		// 对于存在的账号数据叠加
		for i := 0; i < int(count); i++ {
			go func(n int) {
				var buf bytes.Buffer

				// Add(-1)
				defer wg.Done()

				buf.Reset()

				end := int(math.Min(float64((n+1)*BaseLength), float64(repeatLength)))
				start := n * BaseLength
				for j := start; j < end; j++ {
					buf.WriteString(fmt.Sprintf("update item set owner_id = owner_id + %d where serial = %d;", c, items[j].Serial))
				}

				str := buf.String()
				fmt.Println("批量更新 字节长度 item(owner_id): ", len(str), "byte")
				// 批量更新account_common
				err := db.Exec(str).Error
				if err != nil {
					return
				}
			}(i)

			//
			/*buffer.Reset()

			end := int(math.Min(float64((i+1)*BaseLength), float64(repeatLength)))
			start := i * BaseLength
			for j := start; j < end; j++ {
				buffer.WriteString(fmt.Sprintf("update item set owner_id = owner_id + %d where serial = %d;", c, items[j].Serial))
			}

			str := buffer.String()
			fmt.Println("批量更新 字节长度 item(owner_id): ", len(str), "byte")
			// 批量更新account_common
			err := db.Exec(str).Error
			if err != nil {
				return err
			}*/
		}

		wg.Wait()
	}

	return nil
}
