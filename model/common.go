package model

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/config"
	"github.com/panjf2000/ants"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const BaseLength = 1000
const BaseLength1000 = 100
const WanLength = 1000

const (
	Ac = iota + 1
	ItemC
	EquipC
	RoleDataC
	GuildC
	GuildMemberC
	JiMaiItemC
	MailItemC
	FriendC
	BuffC
	EnemyC
	BlackListC
	LeftMsgC
	MapLimitC
	RoleExtendC
	SignC
	SkillC
	TitleC
	SystemIBC
	SystemIIC
	SystemISC
	SystemSBC
	SystemSIC
	SystemSSC
	RewardC
	RoleDataUp1C
	RoleDataUp2C
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
		// 此会话禁用事务
		/*err := db.Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Create(items).Error
		if err != nil {
			return err
		}*/

		/////////////////////////////////////////////////////////////////////////////////////////
		repeatLength := len(items)
		count := math.Ceil(float64(repeatLength) / float64(WanLength))

		var wg sync.WaitGroup
		wg.Add(int(count))
		// 增加池子来多协程一起去插入数据
		// 通过context来控制父子的协程绑定关系
		ctx, cancel := context.WithCancel(context.Background())
		p, _ := ants.NewPoolWithFunc(runtime.NumCPU()-1, func(i interface{}) {
			defer wg.Done()

			tempItems := i.([]*Item)
			err := db.Table("item").WithContext(ctx).Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: WanLength, SkipDefaultTransaction: true}).Create(tempItems).Error
			if err != nil {
				logm.ErrorfE("batch save item 1 [%s]", err.Error())
				cancel()
				return
			}
		})

		defer p.Release()

		for i := 0; i < int(count); i++ {
			end := int(math.Min(float64((i+1)*WanLength), float64(repeatLength)))
			start := i * WanLength
			err := p.Invoke(items[start:end])
			if err != nil {
				logm.ErrorfE("batch save item 2 [%s]", err.Error())
				cancel()
				break
			}
		}

		wg.Wait()

		select {
		// 出现错误直接返回,手动解除
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

	case EquipC:
		equips := arr.([]*Equip)
		// 此会话禁用事务
		err := db.Table("equip").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(equips).Error
		if err != nil {
			return err
		}

		//repeatLength := len(equips)
		//count := math.Ceil(float64(repeatLength) / float64(WanLength))
		//
		//var wg sync.WaitGroup
		//wg.Add(int(count))
		//// 增加池子来多协程一起去插入数据
		//// 通过context来控制父子的协程绑定关系
		//ctx, cancel := context.WithCancel(context.Background())
		//p, _ := ants.NewPoolWithFunc(runtime.NumCPU()-1, func(i interface{}) {
		//	defer wg.Done()
		//
		//	tempItems := i.([]*Equip)
		//	err := db.Table("equip").WithContext(ctx).Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: WanLength, SkipDefaultTransaction: true}).Create(tempItems).Error
		//	if err != nil {
		//		logm.ErrorfE("batch save equip 1 [%s]", err.Error())
		//		cancel()
		//		return
		//	}
		//})
		//
		//for i := 0; i < int(count); i++ {
		//	end := int(math.Min(float64((i+1)*WanLength), float64(repeatLength)))
		//	start := i * WanLength
		//	err := p.Invoke(equips[start:end])
		//	if err != nil {
		//		logm.ErrorfE("batch save equip 2 [%s]", err.Error())
		//		cancel()
		//		break
		//	}
		//}
		//
		//wg.Wait()
		//
		//select {
		//// 出现错误直接返回,手动解除
		//case <-ctx.Done():
		//	return ctx.Err()
		//default:
		//}

	case RoleDataC:
		roles := arr.([]*RoleData)
		// 保存所有的数据到数据库
		err := db.Table("role_data").Table("role_data").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 200, SkipDefaultTransaction: true}).Create(roles).Error
		if err != nil {
			return err
		}
	case GuildC:
		guilds := arr.([]*Guild)
		// save guild from b database to a database
		err := db.Table("guild").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Create(guilds).Error
		if err != nil {
			return err
		}
	case GuildMemberC:
		guildMembers := arr.([]*GuildMember)
		// save guild_member from b database to a database
		err := db.Table("guild_member").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 5000, SkipDefaultTransaction: true}).Create(guildMembers).Error
		if err != nil {
			return err
		}
	case JiMaiItemC:
		jiMais := arr.([]*JimaiItem)
		// save guild_member from b database to a database
		err := db.Table("jimai_item").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(jiMais).Error
		if err != nil {
			return err
		}
	case MailItemC:
		mails := arr.([]*Mail)
		// save mail from b database to a database
		err := db.Table("mail").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 3000, SkipDefaultTransaction: true}).Create(mails).Error
		if err != nil {
			return err
		}
	case FriendC:
		friends := arr.([]*Friend)
		// save friend from b database to a database
		err := db.Table("friend").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 20000, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(friends).Error
		if err != nil {
			return err
		}
	case BuffC:
		buffs := arr.([]*Buff)
		// save friend from b database to a database
		err := db.Table("buff").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 5000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(buffs).Error
		if err != nil {
			return err
		}
	case EnemyC:
		enemys := arr.([]*Enemy)
		// save enemy from b database to a database
		err := db.Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(enemys).Error
		if err != nil {
			return err
		}
	case BlackListC:
		blacks := arr.([]*BlackList)
		// save blacklist from b database to a database
		err := db.Table("blacklist").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 20000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(blacks).Error
		if err != nil {
			return err
		}
	case LeftMsgC:
		msgs := arr.([]*LeftMsg)
		// save left_msg from b database to a database
		err := db.Table("left_msg").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 20000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(msgs).Error
		if err != nil {
			return err
		}
	case MapLimitC:
		limits := arr.([]*MapLimit)
		// save map_limit from b database to a database
		err := db.Table("map_limit").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(limits).Error
		if err != nil {
			return err
		}
	case RoleExtendC:
		extends := arr.([]*RoleExtend)
		// save role_extend from b database to a database
		err := db.Table("role_extend").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 20000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(extends).Error
		if err != nil {
			return err
		}
	case SignC:
		signs := arr.([]*Sign)
		// save sign from b database to a database
		err := db.Table("sign").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(signs).Error
		if err != nil {
			return err
		}
	case SkillC:
		skills := arr.([]*Skill)
		// save skill from b database to a database
		err := db.Table("skill").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(skills).Error
		if err != nil {
			return err
		}
	case TitleC:
		ts := arr.([]*Title)
		// save title from b database to a database
		err := db.Table("title").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(ts).Error
		if err != nil {
			return err
		}
	case SystemIBC:
		ts := arr.([]*SystemIB)
		// save title from b database to a database
		err := db.Table("system_i_b").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(ts).Error
		if err != nil {
			return err
		}
	case SystemIIC:
		s := arr.([]*SystemII)
		err := db.Table("system_i_i").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(s).Error
		if err != nil {
			return err
		}

	case SystemISC:
		s := arr.([]*SystemIS)
		err := db.Table("system_i_s").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(s).Error
		if err != nil {
			return err
		}
	case SystemSBC:
		s := arr.([]*SystemSB)
		err := db.Table("system_s_b").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(s).Error
		if err != nil {
			return err
		}
	case SystemSIC:
		s := arr.([]*SystemSI)
		err := db.Table("system_s_i").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(s).Error
		if err != nil {
			return err
		}
	case SystemSSC:
		s := arr.([]*SystemSS)
		err := db.Table("system_s_s").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000 /*65535/字段个数*/, SkipDefaultTransaction: true}).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(s).Error
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
	case ItemC:
		///////////////////
		//1
		iue := arr.([]*Item)

		//length := len(iue)
		//
		//var tempItem = make([]*Item, length)
		//for k, v := range iue {
		//	tempItem[k] = new(Item)
		//
		//	tempItem[k].Serial = v.Serial
		//	tempItem[k].StrExternData = v.StrExternData
		//}

		//arr = nil

		err := db.Table("item").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "serial"}},
			DoUpdates: clause.AssignmentColumns([]string{"strdwExternData"}),
		}).Create(iue).Error

		if err != nil {
			return err
		}
	/////////////////////////

	////////////////////////
	//2
	//var err1 error
	//for _, v := range iue {
	//	err1 = db.Exec("update item set strdwExternData=? where serial=?;", v.StrExternData, v.Serial).Error
	//	if err1 != nil {
	//		return err1
	//	}
	//}

	////////////////////////
	//3
	//var buffer bytes.Buffer
	//
	//BaseLength := 1000
	//repeatLength := len(iue)
	//count := math.Ceil(float64(repeatLength) / float64(BaseLength))
	//
	//// 对于存在的账号数据叠加
	//for i := 0; i < int(count); i++ {
	//	var tempArr []interface{}
	//	buffer.Reset()
	//
	//	end := int(math.Min(float64((i+1)*BaseLength), float64(repeatLength)))
	//	start := i * BaseLength
	//
	//	for j := start; j < end; j++ {
	//		buffer.WriteString("update item set strdwExternData=? where serial=?;")
	//
	//		tempArr = append(tempArr, iue[j].StrExternData, iue[j].Serial)
	//		//db.Exec("update item set strdwExternData = ? where serial = ?", iue[j].StrExternData, iue[j].Serial)
	//
	//	}
	//
	//	str := buffer.String()
	//
	//	fmt.Println("批量更新 字节长度 item: ", len(str), "byte")
	//
	//	// 批量更新item
	//	err := db.Exec(str, tempArr...).Error
	//	if err != nil {
	//		return err
	//	}
	//}

	case EquipC:
		iue := arr.([]*EquipEx)

		length := len(iue)

		var tempEquip = make([]*Equip, length)
		for k, v := range iue {
			tempEquip[k] = new(Equip)

			tempEquip[k].Serial = v.Serial
			tempEquip[k].EquipExAtt = v.EquipExAtt
			tempEquip[k].EquipAddAtt = v.EquipAddAtt
		}

		arr = nil

		runtime.GC()

		err := db.Table("equip").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "serial"}},
			DoUpdates: clause.AssignmentColumns([]string{"equip_ex_att", "equip_add_att"}),
		}).Create(tempEquip).Error

		if err != nil {
			return err
		}

	case RoleDataUp2C:
		iue := arr.([]*RoleData)

		err := db.Table("role_data").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 200, SkipDefaultTransaction: true}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"script_data", "role_help", "key_info", "strdwExternData", "extern_data2"}),
		}).Create(iue).Error

		if err != nil {
			return err
		}
	case RoleDataUp1C:
		//iue := arr.([]*RoleData)
		//
		//err := db.Table("role_data").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 200, SkipDefaultTransaction: true}).Clauses(clause.OnConflict{
		//	Columns:   []clause.Column{{Name: "role_id"}},
		//	DoUpdates: clause.AssignmentColumns([]string{""}),
		//}).Create(iue).Error
		//
		//if err != nil {
		//	return err
		//}
		rd := arr.([]*RoleData)

		var buffer bytes.Buffer

		repeatLength := len(rd)
		count := math.Ceil(float64(repeatLength) / float64(20000))

		// 对于存在的账号数据叠加
		for i := 0; i < int(count); i++ {
			buffer.Reset()

			end := int(math.Min(float64((i+1)*20000), float64(repeatLength)))
			start := i * 20000
			for j := start; j < end; j++ {
				buffer.WriteString(fmt.Sprintf("update role_data set role_name='%s', rofbid_flag=%d, role_name_origin='%s' where role_id=%d;",
					rd[j].RoleName, rd[j].RofbidFlag, rd[j].RoleNameOrigin, rd[j].RoleID))
			}

			str := buffer.String()
			fmt.Println("批量更新 字节长度 role_data: ", len(str), "byte")
			// 批量更新account_common
			err := db.Exec(str).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type ItemEx struct {
	Serial  int64  `gorm:"primary_key;column:serial"`
	ItenOld int64  `gorm:"column:item_old"`
	Source  string `gorm:"size(100);column:source"`
}

type ItemUpEx struct {
	Serial        int64  `gorm:"primary_key;column:serial"`
	StrExternData []byte `gorm:"column:strdwExternData"`
}

type EquipEx struct {
	Serial      int64  `gorm:"primary_key;column:serial"`
	EquipExAtt  []byte `gorm:"column:equip_ex_att"`
	EquipAddAtt []byte `gorm:"column:equip_add_att"`
}

type MailEx struct {
	MailID uint32 `gorm:"primary_key;column:mail_id"`
}

type RoleDataEx struct {
	RoleID         uint32 `gorm:"primary_key;column:role_id"`
	RoleName       string `gorm:"not null;type:varchar(32);column:role_name"`
	RofbidFlag     int8   `gorm:"not null;column:rofbid_flag"`
	RoleNameOrigin string `gorm:"column:role_name_origin;char(32)"`
}

// 根据时区转换为对应格式的时间字符串
func TimeToStr(t time.Time) string {
	timeStr := t.Format("2006-01-02 15:04:05")                                 //转化所需模板
	loc, _ := time.LoadLocation("Local")                                       //获取时区
	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc) //使用模板在对应时区转化为time.time类型

	return formatTime.Format("2006-01-02 15:04:05")
}

func TimeToTime(t time.Time) time.Time {
	timeStr := t.Format("2006-01-02 15:04:05")                                 //转化所需模板
	loc, _ := time.LoadLocation("Local")                                       //获取时区
	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc) //使用模板在对应时区转化为time.time类型

	return formatTime
}

// 时间戳转字符串
func TimeStampToStr(t int64) string {
	formatTimeStr := time.Unix(t, 0).Format("2006-01-02 15:04:05")
	return formatTimeStr
}

// 根据字符串转换位对应的时间
func StrToTime(str string) (*time.Time, error) {
	loc, _ := time.LoadLocation("Local") //获取时区
	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)

	if err != nil {
		return nil, err
	} else {
		return &formatTime, nil
	}
}

func StrToTimeStamp(str string) (int64, error) {
	loc, _ := time.LoadLocation("Local") //获取时区
	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)

	if err != nil {
		return 0, err
	} else {
		return formatTime.Unix(), nil
	}
}

// err 最大20个字符
// batchid 存时间
// sqlexeeption 20个字符
// errmsg 自定义日志 500个字符
// err				batch id						err msg						sql exeeption time
// item_del10110	0	item_del表数据合并耗时统计10110 merge_log [192.168.1.137 ]	合并成功	2021-01-05 08:56:00
func InsertMergerLog(err, errmsg, sqlexeeption string, batchid int) {
	GDB1.Exec("insert into merger_log(err_log, batchid, errmsg, sqlexeeption, create_time) values(?, ?, ?, ?, ?);", err, batchid, errmsg, sqlexeeption, TimeToTime(time.Now()))
}

// 检测是不是有效的IP
func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

func IsValidIp4(s string) bool {
	_, n := ParseIP(s)

	if n != 4 {
		return false
	}

	return true
}

func MysqlRealEscapeString(value string) string {
	replace := map[string]string{
		"\\":   "\\\\",
		"'":    `\'`,
		"\\0":  "\\\\0",
		"\r":   "\\",
		"\n":   "\\",
		`"`:    `\"`,
		"\x1a": "\\Z",
	}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

func PrintStatus(exMsg string) {
	if !config.Config.OpenMem {
		return
	}

	pid := os.Getpid()

	pro, err := process.NewProcess(int32(pid))
	if err != nil {
		return
	}

	stat, err := pro.MemoryInfo()
	if err != nil {
		return
	}

	str := fmt.Sprintf(exMsg+" 内存指标数据 rss[%d]mb vms[%d]mb", stat.RSS/1024/1024, stat.VMS/1024/1024)

	logm.InfofE(str)
}

func SetLastMergerInfo(ip string, tt int64) {
	LastMergerIp = ip
	LastMergerTime = tt
}
