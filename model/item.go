package model

import "gorm.io/gorm"

/*
  `serial` bigint NOT NULL AUTO_INCREMENT,
  `num` smallint NOT NULL DEFAULT '1' COMMENT '数量',
  `type_id` int unsigned NOT NULL COMMENT '类型id',
  `bind` tinyint NOT NULL COMMENT '是否绑定',
  `lock_state` tinyint NOT NULL DEFAULT '0' COMMENT '物品锁定状态(1:锁定状态)',
  `use_times` int NOT NULL DEFAULT '0' COMMENT '物品使用次数',
  `first_gain_time` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '2000-00-00 00:00:00' COMMENT '第一个玩家获得该道具的时间',
  `create_mode` tinyint NOT NULL COMMENT '创建方式：gm，npc，任务和玩家等',
  `create_id` int unsigned NOT NULL COMMENT '创建方式中对应的id',
  `creator_id` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '创建者id',
  `create_time` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '创建时间',
  `owner_id` int unsigned NOT NULL COMMENT '角色id',
  `account_id` int unsigned NOT NULL COMMENT '帐号id(各角色共享物品用)',
  `container_type_id` tinyint unsigned NOT NULL COMMENT '容器类型',
  `suffix` tinyint unsigned NOT NULL COMMENT '容器中索引',
  `name_id` int unsigned NOT NULL DEFAULT '4294967295' COMMENT '名字id',
  `bind_time` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '2000-00-00 00:00:00' COMMENT '绑定时间',
  `script_data1` int unsigned NOT NULL DEFAULT '0',
  `script_data2` int unsigned NOT NULL DEFAULT '0',
  `create_bind` tinyint NOT NULL DEFAULT '0' COMMENT '创建是否绑定',
  `strdwExternData` blob,
  `item_old` bigint DEFAULT NULL,
  `source` varchar(100) DEFAULT NULL,

*/

type Item struct {
	Serial          int64  `gorm:"primary_key;AUTO_INCREMENT;column:serial"`
	Num             int16  `gorm:"default:1;column:num"`
	TypeID          uint32 `gorm:"column:type_id"`
	Bind            int8   `gorm:"column:bind"`
	LockState       int8   `gorm:"column:lock_state"`
	UseTimes        int32  `gorm:"column:use_times"`
	FirstGainTime   string `gorm:"type:char(20);column:first_gain_time"`
	CreateNode      int8   `gorm:"column:create_mode"`
	CreateID        uint32 `gorm:"column:create_id"`
	CreatorID       uint32 `gorm:"column:creator_id"`
	CreateTime      string `gorm:"char(20);column:create_time"`
	OwnerID         uint32 `gorm:"column:owner_id"`
	AccountID       uint32 `gorm:"column:account_id:"`
	ContainerTypeID uint8  `gorm:"column:container_type_id"`
	Suffix          uint8  `gorm:"column:suffix"`
	NameID          uint32 `gorm:"default:4294967295;column:name_id"`
	BindTime        string `gorm:"char(20);column:bind_time"`
	ScriptData1     uint32 `gorm:"column:script_data1"`
	ScriptData2     uint32 `gorm:"column:script_data2"`
	CreateBind      int8   `gorm:"column:create_bind"`
	StrExternData   []byte `gorm:"column:strdwExternData"`
	ItenOld         int64  `gorm:"column:item_old"`
	Source          string `gorm:"size(100);column:source"`
}

func (m *Item) TableName() string {
	return "item"
}

func (m *Item) Query(gdb *gorm.DB, where ...interface{}) []Item {
	var items []Item
	gdb.Find(&items, where)
	return items
}

func (m *Item) Read(gdb *gorm.DB, fields ...string) (items []Item, err error) {
	if len(fields) == 0 {
		err = gdb.Find(&items).Error
		return
	}

	err = gdb.Select(fields).Find(&items).Error
	return
}

func (m *Item) GetCount(gdb *gorm.DB) (count int64) {
	gdb.Table(m.TableName()).Count(&count)
	return
}
