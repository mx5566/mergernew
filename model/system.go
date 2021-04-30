package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
	"time"
)

// system_i_b system_i_i system_i_s
// system_s_b system_s_i system_s_s
// 6 table struct
type SystemIB struct {
	SysKey       int64  `gorm:"primaryKey;column:sys_key;autoIncrement:false"`
	SysValue     []byte `gorm:"column:sys_value"`
	SysClearType int16  `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16  `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemIB(db1, db2 *gorm.DB) error {
	var sib1 []*SystemIB
	var sib2 []*SystemIB

	err := db1.Table("system_i_b").Find(&sib1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_i_b").Find(&sib2).Error
	if err != nil {
		return err
	}

	// logic handle
	for i := 0; i < len(sib1); i++ {
		if sib1[i].SysAreaType == 1 {
			sib1 = append(sib1[0:i], sib1[i+1:]...)
			i--
		} else if sib1[i].SysAreaType == 6 {
			sib1 = append(sib1[0:i], sib1[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(sib2); i++ {
		sat := sib2[i].SysAreaType
		if sat != 1 {
			sib2 = append(sib2[0:i], sib2[i+1:]...)
			i--
		}
	}

	sib1 = append(sib1, sib2...)
	err = db1.Exec("truncate table system_i_b;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemIBC, sib1)
	if err != nil {
		return err
	}

	return nil
}

type SystemII struct {
	SysKey       int64 `gorm:"primaryKey;column:sys_key;autoIncrement:false"`
	SysValue     int64 `gorm:"column:sys_value;default:0"`
	SysClearType int16 `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16 `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemII(db1, db2 *gorm.DB) error {
	var sii1 []*SystemII
	var sii2 []*SystemII

	err := db1.Table("system_i_i").Find(&sii1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_i_i").Find(&sii2).Error
	if err != nil {
		return err
	}

	// k1 SysKey k2 SysAreaType v2 sii1 / sii2 的数组index
	mapTypeV := make(map[int64]map[int16]int)
	// logic handle
	for i := 0; i < len(sii1); i++ {
		if sii1[i].SysAreaType == 1 {
			sii1 = append(sii1[0:i], sii1[i+1:]...)
			i--
			continue
		} else if sii1[i].SysAreaType == 6 {
			sii1 = append(sii1[0:i], sii1[i+1:]...)
			i--
			continue
		}
		sat := sii1[i].SysAreaType
		skey := sii1[i].SysKey

		if mapTypeV[skey] == nil {
			mapTypeV[skey] = make(map[int16]int)
		}

		mapTypeV[skey][sat] = i
	}

	deleteIndex := []int{}
	for i := 0; i < len(sii2); i++ {
		sat := sii2[i].SysAreaType
		skey := sii2[i].SysKey

		if sat == 6 {
			sii2 = append(sii2[0:i], sii2[i+1:]...)
			i--
			continue
		} else if sat == 1 {
			continue
		} else if sat == 2 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[2]; ok {
					if sii1[v2].SysValue > sii2[i].SysValue {
						sii2 = append(sii2[0:i], sii2[i+1:]...)
						i--
					} else {
						//sii1 = append(sii1[0:v2], sii1[v2+1:]...)
						deleteIndex = append(deleteIndex, v2)
					}
				}
			}
			continue
		} else if sat == 3 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[3]; ok {
					if sii1[v2].SysValue < sii2[i].SysValue {
						sii2 = append(sii2[0:i], sii2[i+1:]...)
						i--
					} else {
						deleteIndex = append(deleteIndex, v2)
						//sii1 = append(sii1[0:v2], sii1[v2+1:]...)
					}
				}
			}
			continue
		} else if sat == 4 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[4]; ok {
					sii1[v2].SysValue += sii2[i].SysValue
					sii2 = append(sii2[0:i], sii2[i+1:]...)
					i--
				}
			}
			continue
		} else {
			sii2 = append(sii2[0:i], sii2[i+1:]...)
			i--
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		sii1 = append(sii1[0:value-count], sii1[value+1-count:]...)
		count++
	}

	sii1 = append(sii1, sii2...)
	err = db1.Exec("truncate table system_i_i;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemIIC, sii1)
	if err != nil {
		return err
	}

	return nil
}

type SystemIS struct {
	SysKey       int64  `gorm:"primaryKey;column:sys_key;autoIncrement:false"`
	SysValue     string `gorm:"column:sys_value;char(200)"`
	SysClearType int16  `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16  `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemIS(db1, db2 *gorm.DB) error {
	var sis1 []*SystemIS
	var sis2 []*SystemIS

	err := db1.Table("system_i_s").Find(&sis1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_i_s").Find(&sis2).Error
	if err != nil {
		return err
	}

	// k1 SysKey k2 SysAreaType v2 sis1 / sis2 的数组index
	mapTypeV := make(map[int64]map[int16]int)
	// logic handle
	for i := 0; i < len(sis1); i++ {
		if sis1[i].SysAreaType == 1 {
			sis1 = append(sis1[0:i], sis1[i+1:]...)
			i--
			continue
		} else if sis1[i].SysAreaType == 6 {
			sis1 = append(sis1[0:i], sis1[i+1:]...)
			i--
			continue
		}
		sat := sis1[i].SysAreaType
		skey := sis1[i].SysKey

		if mapTypeV[skey] == nil {
			mapTypeV[skey] = make(map[int16]int)
		}

		mapTypeV[skey][sat] = i
	}

	for i := 0; i < len(sis2); i++ {
		sat := sis2[i].SysAreaType
		skey := sis2[i].SysKey

		if sat == 6 {
			sis2 = append(sis2[0:i], sis2[i+1:]...)
			i--
			continue
		} else if sat == 1 {
			continue
		} else if sat == 5 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[5]; ok {
					sis1[v2].SysValue += sis2[i].SysValue
					sis2 = append(sis2[0:i], sis2[i+1:]...)
					i--
				}
			}
			continue
		} else {
			sis2 = append(sis2[0:i], sis2[i+1:]...)
			i--
		}
	}

	sis1 = append(sis1, sis2...)
	err = db1.Exec("truncate table system_i_s;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemISC, sis1)
	if err != nil {
		return err
	}

	return nil
}

type SystemSB struct {
	SysKey       string `gorm:"primaryKey;column:sys_key;char(200)"`
	SysValue     []byte `gorm:"column:sys_value"`
	SysClearType int16  `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16  `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemSB(db1, db2 *gorm.DB) error {
	var ssb1 []*SystemSB
	var ssb2 []*SystemSB

	err := db1.Table("system_s_b").Find(&ssb1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_s_b").Find(&ssb2).Error
	if err != nil {
		return err
	}

	// logic handle
	// k1 SysKey k2 SysAreaType v2 sss1 / sss2 的数组index
	mapTypeV := make(map[string]map[int16]int)
	// logic handle
	for i := 0; i < len(ssb1); i++ {
		if ssb1[i].SysAreaType == 1 {
			ssb1 = append(ssb1[0:i], ssb1[i+1:]...)
			i--
			continue
		} else if ssb1[i].SysAreaType == 6 {
			ssb1 = append(ssb1[0:i], ssb1[i+1:]...)
			i--
			continue
		}
		sat := ssb1[i].SysAreaType
		skey := ssb1[i].SysKey

		if mapTypeV[skey] == nil {
			mapTypeV[skey] = make(map[int16]int)
		}

		mapTypeV[skey][sat] = i
	}

	for i := 0; i < len(ssb2); i++ {
		sat := ssb2[i].SysAreaType
		skey := ssb2[i].SysKey

		if sat == 6 {
			ssb2 = append(ssb2[0:i], ssb2[i+1:]...)
			i--
			continue
		} else if sat == 1 {
			continue
		} else if sat == 5 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[5]; ok {
					ssb1[v2].SysValue = append(ssb1[v2].SysValue, ssb2[i].SysValue...)
					ssb2 = append(ssb2[0:i], ssb2[i+1:]...)
					i--
				}
			}
			continue
		} else {
			ssb2 = append(ssb2[0:i], ssb2[i+1:]...)
			i--
		}
	}

	ssb1 = append(ssb1, ssb2...)
	err = db1.Exec("truncate table system_s_b;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemSBC, ssb1)
	if err != nil {
		return err
	}

	return nil
}

type SystemSI struct {
	SysKey       string `gorm:"primaryKey;column:sys_key;char(200)"`
	SysValue     int64  `gorm:"column:sys_value;default:0"`
	SysClearType int16  `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16  `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemSI(db1, db2 *gorm.DB) error {
	var ssi1 []*SystemSI
	var ssi2 []*SystemSI

	err := db1.Table("system_s_i").Find(&ssi1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_s_i").Find(&ssi2).Error
	if err != nil {
		return err
	}

	// k1 SysKey k2 SysAreaType v2 ssi1 / ssi2 的index
	mapTypeV := make(map[string]map[int16]int)
	// logic handle
	for i := 0; i < len(ssi1); i++ {
		if ssi1[i].SysAreaType == 1 {
			ssi1 = append(ssi1[0:i], ssi1[i+1:]...)
			i--
			continue
		} else if ssi1[i].SysAreaType == 6 {
			ssi1 = append(ssi1[0:i], ssi1[i+1:]...)
			i--
			continue
		}
		sat := ssi1[i].SysAreaType
		skey := ssi1[i].SysKey

		if mapTypeV[skey] == nil {
			mapTypeV[skey] = make(map[int16]int)
		}

		mapTypeV[skey][sat] = i
	}

	deleteIndex := []int{}
	for i := 0; i < len(ssi2); i++ {
		sat := ssi2[i].SysAreaType
		skey := ssi2[i].SysKey

		if sat == 6 {
			ssi2 = append(ssi2[0:i], ssi2[i+1:]...)
			i--
			continue
		} else if sat == 1 {
			continue
		} else if sat == 2 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[2]; ok {
					if ssi1[v2].SysValue > ssi2[i].SysValue {
						ssi2 = append(ssi2[0:i], ssi2[i+1:]...)
						i--
					} else {

						deleteIndex = append(deleteIndex, v2)
						//ssi1 = append(ssi1[0:v2], ssi1[v2+1:]...)
					}
				}
			}
			continue
		} else if sat == 3 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[3]; ok {
					if ssi1[v2].SysValue < ssi2[i].SysValue {
						ssi2 = append(ssi2[0:i], ssi2[i+1:]...)
						i--
					} else {
						deleteIndex = append(deleteIndex, v2)
						//ssi1 = append(ssi1[0:v2], ssi1[v2+1:]...)
					}
				}
			}
			continue
		} else if sat == 4 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[4]; ok {
					ssi1[v2].SysValue += ssi2[i].SysValue
					ssi2 = append(ssi2[0:i], ssi2[i+1:]...)
					i--
				}
			}
			continue
		} else {
			ssi2 = append(ssi2[0:i], ssi2[i+1:]...)
			i--
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		ssi1 = append(ssi1[0:value-count], ssi1[value+1-count:]...)
		count++
	}

	ssi1 = append(ssi1, ssi2...)
	err = db1.Exec("truncate table system_s_i;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemSIC, ssi1)
	if err != nil {
		return err
	}

	return nil
}

type SystemSS struct {
	SysKey       string `gorm:"primaryKey;column:sys_key;char(200)"`
	SysValue     string `gorm:"column:sys_value;char(200)"`
	SysClearType int16  `gorm:"column:sys_clear_type;default:0"`
	SysAreaType  int16  `gorm:"column:sys_area_type;default:0"`
}

func HandleSystemSS(db1, db2 *gorm.DB) error {
	var sss1 []*SystemSS
	var sss2 []*SystemSS

	err := db1.Table("system_s_s").Find(&sss1).Error
	if err != nil {
		return err
	}

	err = db2.Table("system_s_s").Find(&sss2).Error
	if err != nil {
		return err
	}

	// logic handle
	// k1 SysKey k2 SysAreaType v2 sss1 / sss2 的数组index
	mapTypeV := make(map[string]map[int16]int)
	// logic handle
	for i := 0; i < len(sss1); i++ {
		if sss1[i].SysAreaType == 1 {
			sss1 = append(sss1[0:i], sss1[i+1:]...)
			i--
			continue
		} else if sss1[i].SysAreaType == 6 {
			sss1 = append(sss1[0:i], sss1[i+1:]...)
			i--
			continue
		}
		sat := sss1[i].SysAreaType
		skey := sss1[i].SysKey

		if mapTypeV[skey] == nil {
			mapTypeV[skey] = make(map[int16]int)
		}

		mapTypeV[skey][sat] = i
	}

	for i := 0; i < len(sss2); i++ {
		sat := sss2[i].SysAreaType
		skey := sss2[i].SysKey

		if sat == 6 {
			sss2 = append(sss2[0:i], sss2[i+1:]...)
			i--
			continue
		} else if sat == 1 {
			continue
		} else if sat == 5 {
			if v1, ok := mapTypeV[skey]; ok {
				// 判断类型是不是存在
				if v2, ok := v1[5]; ok {
					sss1[v2].SysValue += sss2[i].SysValue
					sss2 = append(sss2[0:i], sss2[i+1:]...)
					i--
				}
			}
			continue
		} else {
			sss2 = append(sss2[0:i], sss2[i+1:]...)
			i--
		}
	}

	sss1 = append(sss1, sss2...)
	err = db1.Exec("truncate table system_s_s;").Error
	if err != nil {
		return err
	}

	//Clauses(clause.Insert{Modifier: "IGNORE"}).
	err = BatchSave(db1, SystemSSC, sss1)
	if err != nil {
		return err
	}

	return nil
}

func HandleSystem(db1, db2 *gorm.DB) error {
	t1 := time.Now().Unix()
	err := HandleSystemIB(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemIB ok ... time[%d]", time.Now().Unix()-t1)

	t1 = time.Now().Unix()
	err = HandleSystemII(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemII ok ... time[%d]", time.Now().Unix()-t1)

	t1 = time.Now().Unix()
	err = HandleSystemIS(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemIS ok ... time[%d]", time.Now().Unix()-t1)

	t1 = time.Now().Unix()
	err = HandleSystemSB(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemSB ok ... time[%d]", time.Now().Unix()-t1)

	t1 = time.Now().Unix()
	err = HandleSystemSI(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemSI ok ... time[%d]", time.Now().Unix()-t1)

	t1 = time.Now().Unix()
	err = HandleSystemSS(db1, db2)
	if err != nil {
		return err
	}
	logm.DebugfE("HandleSystemSS ok ... time[%d]", time.Now().Unix()-t1)

	return nil
}
