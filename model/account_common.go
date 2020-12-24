package model

/*
  `account_id` int unsigned NOT NULL COMMENT '账号id',
  `account_name` char(36) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '账号名',
  `safecode_crc` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '安全码crc32',
  `reset_time` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '重置时间',
  `bag_password_crc` int unsigned NOT NULL DEFAULT '111111' COMMENT '背包密码crc值',
  `baibao_yuanbao` int NOT NULL DEFAULT '0' COMMENT '百宝中元宝',
  `ware_size` smallint NOT NULL DEFAULT '8' COMMENT '仓库大小,默认值为8',
  `ware_silver` int NOT NULL DEFAULT '0' COMMENT '仓库中金钱数',
  `warestep` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '扩充仓库阶段',
  `yuanbao_recharge` int unsigned NOT NULL DEFAULT '0' COMMENT '添加元宝',
  `IsReceive` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '添加奖励领取标志',
  `total_recharge` int unsigned NOT NULL DEFAULT '0' COMMENT '添加累计充值数量',
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
	BaiBaoYuanBao   int32  `gorm:"default:0;not null;column:baibao_yuanbao"`
	WareSize        int16  `gorm:"not null;default 8;column:ware_size"`
	WareSilver      int32  `gorm:"not null;default 0;column:ware_silver"`
	WareStep        uint8  `gorm:"not null;default 1;column:warestep"`
	YuanBaoRecharge uint32 `gorm:"not null;default 0;column:yuanbao_recharge"`
	IsReceive       uint8  `gorm:"not null;default 0;column:IsReceive"`
	TotalRecharge   uint32 `gorm:"not null;default 0;column:total_recharge"`
	ReceiveType     uint16 `gorm:"not null;default 0;column:receive_type"`
	ReceiveTypeEx   uint32 `gorm:"not null;default 0;column:receive_type_ex"`
	WebType         uint32 `gorm:"not null;default 0;column:web_type"`
	Score           uint32 `gorm:"not null;default 0;column:score"`
	Robort          uint8  `gorm:"not null;default 0;column:robort"`
	Forbid_flag     int8   `gorm:"not null;default 0;column:forbid_flag"`
	DataEx          string `gorm:"column:data_ex"`
}

type DataEx struct {
	TodayRecharge    int32  `gorm:"default 0;column:today_recharge"`
	TodayRechargeDay uint32 `gorm:"default 0;column:today_recharge_day"`
}

func (m *AccountCommon) TableName() string {
	return "account_common"
}
