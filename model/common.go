package model

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/panjf2000/ants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"runtime"
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
		err := db.Table("equip").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 1000, SkipDefaultTransaction: true}).Create(equips).Error
		if err != nil {
			return err
		}
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
		err := db.Table("jimai_item").Session(&gorm.Session{PrepareStmt: true, CreateBatchSize: 10000, SkipDefaultTransaction: true}).Create(jiMais).Error
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

	}

	return nil
}

type ItemEx struct {
	Serial  int64  `gorm:"primary_key;column:serial"`
	ItenOld int64  `gorm:"column:item_old"`
	Source  string `gorm:"size(100);column:source"`
}

type MailEx struct {
	MailID uint32 `gorm:"primary_key;column:mail_id"`
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
