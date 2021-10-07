package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

/*
  `account_id` int unsigned NOT NULL COMMENT '账号id',
  `account_name` char(36) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '账号名',
  `safecode_crc` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '安全码crc32',
  `reset_time` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '重置时间',
  `bag_password_crc` int unsigned NOT NULL DEFAULT '111111' COMMENT '背包密码crc值',
  `baibao_yuanbao` int NOT NULL DEFAULT '0' COMMENT '百宝中元宝', 变为 bigint
  `ware_size` smallint NOT NULL DEFAULT '8' COMMENT '仓库大小,默认值为8',
  `ware_silver` int NOT NULL DEFAULT '0' COMMENT '仓库中金钱数',
  `warestep` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '扩充仓库阶段',
  `yuanbao_recharge` int unsigned NOT NULL DEFAULT '0' COMMENT '添加元宝',
  `IsReceive` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '添加奖励领取标志',
  `total_recharge` int unsigned NOT NULL DEFAULT '0' COMMENT '添加累计充值数量', 变为 unsigned  bigint
  `receive_type` smallint unsigned NOT NULL DEFAULT '0' COMMENT '领奖类型',
  `receive_type_ex` int unsigned NOT NULL DEFAULT '0' COMMENT '添加领奖类型扩充',
  `web_type` int unsigned NOT NULL DEFAULT '0' COMMENT '添加网页端领奖类型',
  `score` int unsigned NOT NULL DEFAULT '0' COMMENT '添加玩家积分',
  `robort` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否是机器人',
  `forbid_flag` tinyint NOT NULL DEFAULT '0',
  `data_ex` json DEFAULT NULL,
*/

type AccountCommon struct {
	AccountID       uint32 `gorm:"primary_key;column:account_id"`
	AccountName     string `gorm:"char(36);column:account_name"`
	SafeCodeCrc     uint32 `gorm:"default:4294967295;column:account_name"`
	ResetTime       uint32 `gorm:"default:4294967295;column:reset_time;not null"`
	BagPasswordCrc  uint32 `gorm:"not null; default 111111;column:bag_password_crc"`
	BaiBaoYuanBao   int64  `gorm:"default:0;not null;column:baibao_yuanbao"`
	WareSize        int16  `gorm:"not null;default 8;column:ware_size"`
	WareSilver      int32  `gorm:"not null;default 0;column:ware_silver"`
	WareStep        uint8  `gorm:"not null;default 1;column:warestep"`
	YuanBaoRecharge uint32 `gorm:"not null;default 0;column:yuanbao_recharge"`
	IsReceive       uint8  `gorm:"not null;default 0;column:IsReceive"`
	TotalRecharge   uint64 `gorm:"not null;default 0;column:total_recharge"`
	ReceiveType     uint16 `gorm:"not null;default 0;column:receive_type"`
	ReceiveTypeEx   uint32 `gorm:"not null;default 0;column:receive_type_ex"`
	WebType         uint32 `gorm:"not null;default 0;column:web_type"`
	Score           uint32 `gorm:"not null;default 0;column:score"`
	Robort          uint8  `gorm:"not null;default 0;column:robort"`
	Forbid_flag     int8   `gorm:"not null;default 0;column:forbid_flag"`
	DataEx          string `gorm:"column:data_ex"`
}

type DataEx struct {
	TodayRecharge    int32  `gorm:"default 0;column:today_recharge" json:"today_recharge"`
	TodayRechargeDay uint32 `gorm:"default 0;column:today_recharge_day" json:"today_recharge_day"`
}

func (m *AccountCommon) TableName() string {
	return "account_common"
}

func HandleAccountCommon(db1, db2 *gorm.DB) error {
	// 处理account_common表
	var accounts1 []*AccountCommon
	err := db1.Select("account_id, baibao_yuanbao, total_recharge, yuanbao_recharge, data_ex").Find(&accounts1).Error
	//
	if err != nil {
		return err
	}

	var accounts2 []*AccountCommon
	err = db2.Find(&accounts2).Error
	//
	if err != nil {
		return err
	}

	// 把相同account_id的账号的数据充值元宝合并 data_ex 是个json串需要单独处理
	a2 := make(map[uint32]*AccountCommon)
	a1 := make(map[uint32]*AccountCommon)

	alreadyAccount := make([]*AccountCommon, 0)
	notExistAccount := make([]*AccountCommon, 0)

	for _, v2 := range accounts2 {
		a2[v2.AccountID] = v2
	}

	for _, v1 := range accounts1 {
		a1[v1.AccountID] = v1
	}

	// A为基准
	s := time.Now().Unix()
	day := s / (60 * 60 * 24)
	fmt.Println("+++++++++++++++++++++++", day)
	for _, v2 := range accounts2 {
		if _, ok := a1[v2.AccountID]; !ok {
			// 没有找到对应的账号，直接插入
			if v2.DataEx == "" {
				v2.DataEx = "{}"
			}
			notExistAccount = append(notExistAccount, v2)
			continue
		}

		data := a1[v2.AccountID]

		// 找到对应的数据，数据叠加
		// 溢出处理逻辑
		r1 := uint64(data.BaiBaoYuanBao)
		r2 := uint64(v2.BaiBaoYuanBao)

		if r1+r2 > uint64(math.MaxInt64) {
			v2.BaiBaoYuanBao = math.MaxInt64
		} else {
			v2.BaiBaoYuanBao += data.BaiBaoYuanBao
		}

		//v2.BaiBaoYuanBao += data.BaiBaoYuanBao
		v2.TotalRecharge += data.TotalRecharge
		v2.YuanBaoRecharge += data.YuanBaoRecharge

		// 修改data_ex字段
		d1 := v2.DataEx
		d2 := data.DataEx
		de1 := new(DataEx)
		de2 := new(DataEx)

		if d1 == "" {
			d1 = "{}"
		}

		if d2 == "" {
			d2 = "{}"
		}

		err := json.Unmarshal([]byte(d1), de1)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(d2), de2)
		if err != nil {
			return err
		}

		if day == int64(de1.TodayRechargeDay) && day == int64(de2.TodayRechargeDay) {
			de1.TodayRecharge += de2.TodayRecharge
		} else if day == int64(de2.TodayRechargeDay) {
			de1.TodayRecharge = de2.TodayRecharge
			de1.TodayRechargeDay = de2.TodayRechargeDay
		}

		if de1.TodayRechargeDay != 0 {
			by, err := json.Marshal(&de1)
			if err != nil {
				panic(err)
			}
			v2.DataEx = string(by)
		} else {
			v2.DataEx = "{}"
		}

		alreadyAccount = append(alreadyAccount, v2)

	}

	repeatLength := len(alreadyAccount)
	count := math.Ceil(float64(repeatLength) / float64(BaseLength))

	fmt.Println(count)

	// 对于存在的账号数据叠加
	err = BatchUpdate(GDB1, Ac, alreadyAccount)
	if err != nil {
		return err
	}

	// 不存在的数据批量插入
	err = BatchSave(GDB1, Ac, notExistAccount)
	if err != nil {
		return err
	}
	return nil
}
