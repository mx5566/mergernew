package model

import (
	"database/sql"
	"github.com/mx5566/logm"
	"gorm.io/gorm"
	"time"
)

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
	Serial          int64          `gorm:"primary_key;autoIncrement;column:serial"`
	Num             int16          `gorm:"default:1;column:num;not null"`
	TypeID          uint32         `gorm:"column:type_id;not null"`
	Bind            int8           `gorm:"column:bind;not null"`
	LockState       int8           `gorm:"column:lock_state;not null"`
	UseTimes        int32          `gorm:"column:use_times;not null"`
	FirstGainTime   string         `gorm:"type:char(20);column:first_gain_time;not null"`
	CreateNode      int8           `gorm:"column:create_mode;not null"`
	CreateID        uint32         `gorm:"column:create_id;not null"`
	CreatorID       uint32         `gorm:"column:creator_id;not null"`
	CreateTime      string         `gorm:"char(20);column:create_time;not null"`
	OwnerID         uint32         `gorm:"column:owner_id;not null"`
	AccountID       uint32         `gorm:"column:account_id;not null"`
	ContainerTypeID uint8          `gorm:"column:container_type_id;not null"`
	Suffix          uint16         `gorm:"column:suffix;not null"`
	NameID          uint32         `gorm:"default:4294967295;column:name_id;not null"`
	BindTime        string         `gorm:"char(20);column:bind_time;not null"`
	ScriptData1     uint32         `gorm:"column:script_data1;not null"`
	ScriptData2     uint32         `gorm:"column:script_data2;not null"`
	CreateBind      int8           `gorm:"column:create_bind;not null"`
	StrExternData   []byte         `gorm:"column:strdwExternData"`
	ItenOld         sql.NullInt64  `gorm:"column:item_old"`
	Source          sql.NullString `gorm:"size(100);column:source"`
}

func (m *Item) TableName() string {
	return "item"
}

func HandleItemRelation(db1, db2 *gorm.DB, increaseNum uint32, mapMailIDs map[uint32]uint32) error {
	//err := db1.Exec("ALTER TABLE item MODIFY serial BIGINT(21)  AUTO_INCREMENT;").Error
	//if err != nil {
	//	return err
	//}

	var MAXItemID int64
	type MaxStruct struct {
		Max int64 `json:"max"`
		Min int64 `json:"min"`
	}

	mapItems := make(map[int64]*ItemEx)

	var result MaxStruct
	err := db1.Raw("select max(serial) as max from item;").Scan(&result).Error
	if err != nil {
		return err
	}

	MAXItemID = result.Max
	// 表示1号数据库没有物品数据
	startID := int64(500000000000)
	if MAXItemID == 0 {
		startID = 499999999999
	} else {
		startID = MAXItemID
	}

	t1 := time.Now().Unix()
	t2 := time.Now().Unix()
	delta := t2 - t1

	SetCurrent(Stage_31, "", delta)
	var items []*Item
	err = db2.Select("serial, num, type_id, bind, lock_state, use_times, first_gain_time, create_mode, create_id, creator_id, create_time, owner_id, account_id, container_type_id, suffix, name_id, bind_time, script_data1, script_data2, create_bind, strdwExternData, item_old, source").Find(&items).Error
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - t1
	logm.DebugfE("handle select item delta[%d] start-item-id[%d]", delta, startID)

	// 遍历所有的物品
	for index, value := range items {
		items[index].ItenOld.Int64 = value.Serial
		items[index].ItenOld.Valid = true
		items[index].Source.String = "_copy"
		items[index].Source.Valid = true

		// 把非邮件的物品所有者加指定值
		if value.ContainerTypeID != 11 {
			items[index].OwnerID += increaseNum
		} else {
			if v, ok := mapMailIDs[items[index].OwnerID]; ok {
				items[index].OwnerID = v
			}
			//items[index].OwnerID += increaseNumMail
		}

		// 把所有8的 好像是礼包
		if value.CreateNode == 8 {
			items[index].CreateID += increaseNum
		}

		vv := new(ItemEx)

		startID++

		vv.Serial = startID
		vv.ItenOld = value.Serial
		vv.Source = "_copy"
		mapItems[vv.ItenOld] = vv

		items[index].Serial = startID
	}
	t2 = time.Now().Unix()
	delta = t2 - t1
	logm.DebugfE("handle item memory handle ... time[%d]", delta)

	SetCurrent(Stage_40, "", delta)
	// handle equip
	err = HandleEquip(db1, db2, mapItems)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle equip batch save ... time[%d]", delta)

	SetCurrent(Stage_50, "", delta)
	//
	// handle jimai_item
	err = HandleJiMai(db1, db2, mapItems)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle jimai batch save ... time[%d]", delta)

	SetCurrent(Stage_60, "", delta)
	//
	// handle mail item
	err = HandleMailItem(db1, db2, mapItems, mapMailIDs)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle mail item batch save ... time[%d]", delta)

	SetCurrent(Stage_70, "", delta)
	for k1, _ := range items {
		// 对于null值的处理方式
		items[k1].ItenOld.Valid = false
		items[k1].ItenOld.Int64 = 0
		items[k1].Source.Valid = false
		items[k1].Source.String = ""
	}

	dropTime := time.Now().Unix()
	// 删除item的所有索引
	err = db1.Exec("drop index owner_id on item;drop index account_id on item;drop index container_type_id on item;").Error
	if err != nil {
		return err
	}

	t2 = time.Now().Unix()
	delta = t2 - dropTime
	logm.DebugfE("drop item index  ... time[%d]", delta)

	// 保存之前需要把所有的item_old 和source字段设置为空
	// 批量插入被合数据库表item
	err = BatchSave(db1, ItemC, items)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle item save ... time[%d] start[%s]", delta, TimeToStr(time.Now()))

	// add index item
	err = db1.Exec("CREATE INDEX  owner_id on item (owner_id); CREATE INDEX  account_id on item (account_id, container_type_id);CREATE INDEX  container_type_id on item (container_type_id);").Error
	if err != nil {
		return err
	}
	logm.DebugfE("create index item ... time[%d] ", time.Now().Unix()-t2)

	return nil
}
