package model

import (
	"github.com/mx5566/logm"
	"gorm.io/gorm"
	"time"
)

// guild struct
type Guild struct {
	ID                uint64    `gorm:"primary_key;column:id"`
	Name              []byte    `gorm:"column:name"`
	CreaterNameID     uint32    `gorm:"not null;column:creater_name_id"`
	LeaderID          uint32    `gorm:"not null;column:leader_id"`
	SpecialState      uint32    `gorm:"not null;column:special_state"`
	Level             int8      `gorm:"not null;column:level"`
	HoldCity0         uint8     `gorm:"not null;column:hold_city0"`
	HoldCity1         uint8     `gorm:"not null;column:hold_city1"`
	HoldCity2         uint8     `gorm:"not null;column:hold_city2"`
	Fund              int32     `gorm:"not null;column:fund"`
	Material          int32     `gorm:"not null;column:material"`
	Reputation        int32     `gorm:"not null;column:reputation"`
	DailyCost         int32     `gorm:"not null;column:daily_cost"`
	Peace             int16     `gorm:"not null;column:peace"`
	Tenet             []byte    `gorm:"column:tenet"`
	Symbol            []byte    `gorm:"column:symbol"`
	CreateTime        time.Time `gorm:"not null;column:create_time"`
	Formal            uint8     `gorm:"not null;column:formal"`
	SignNum           uint8     `gorm:"not null;column:sign_num"`
	SignData          []byte    `gorm:"not null;column:sign_data"`
	UpLevelTime       string    `gorm:"char(20);not null;column:uplevel_time"`
	EnemyData         []byte    `gorm:"column:enemy_data"`
	LeagueID          uint32    `gorm:"not null;column:league_id;default:4294967295"`
	UnLeagueTime      string    `gorm:"char(20);not null;column:unleague_time"`
	Prosperity        uint32    `gorm:"not null;column:prosperity"`
	PosName           []byte    `gorm:"column:pos_name"`
	PosPower          []byte    `gorm:"column:pos_power"`
	Dkp               []byte    `gorm:"column:dkp"`
	ChangeDkp         uint8     `gorm:"not null;column:change_dkp"`
	ChangeGuildSymbol string    `gorm:"not null;column:change_guild_symbol;char(20)"`
	SymbolValue       uint32    `gorm:"not null;column:symbol_value"`
	GroupPurchase     uint32    `gorm:"not null;column:group_purchase"`
	RemainSpreadTimes uint32    `gorm:"not null;column:remain_spread_times"`
	Commendation      int8      `gorm:"not null;column:commendation"`
	ScriptData        []byte    `gorm:"column:script_data"`
	Text              []byte    `gorm:"column:text"`
	FamilyName        []byte    `gorm:"column:family_name"`
	NpcName           []byte    `gorm:"column:npc_name"`
	ChangeName        uint8     `gorm:"not null;column:change_name"`
	DaoGao            []byte    `gorm:"column:daogao"`
	JujueTime         string    `gorm:"column:jujue_time;char(20)"`
	RefuseApply       int16     `gorm:"not null;column:refuse_apply"`
	AttackCity        int32     `gorm:"column:attack_city"`
	GuildRank         int16     `gorm:"not null;column:gulid_rank"`
}

// guild member struct
type GuildMember struct {
	RoleID            uint32    `gorm:"primary_key;column:role_id"`
	GuildID           uint32    `gorm:"not null;column:guild_id"`
	GuildPos          uint8     `gorm:"not null;column:guild_pos"`
	TotalContribution int32     `gorm:"not null;column:total_contribution"`
	CurContribution   int32     `gorm:"not null;column:cur_contribution"`
	Exploit           int32     `gorm:"not null;column:exploit"`
	Salary            int32     `gorm:"not null;column:salary"`
	CanUseGuildWare   int8      `gorm:"not null;column:can_use_guild_ware"`
	Ballot            uint8     `gorm:"not null;column:ballot"`
	Dkp               uint32    `gorm:"not null;column:dkp"`
	JoinTime          time.Time `gorm:"not null;column:join_time"`
	TotalFund         int32     `gorm:"column:total_fund"`
}

// 和角色数据一起处理
func HandleGuildRelation(db1, db2 *gorm.DB, roles1 []*RoleData, lengthRole /*这个是roles1合并数据之前的长度*/ int, mapRoleIds map[uint32]uint32) error {
	guilds1 := make([]*Guild, 0)
	guilds2 := make([]*Guild, 0)

	// load guild
	err := db1.Table("guild").Order("id asc").Find(&guilds1).Error
	if err != nil {
		return err
	}

	// load guild
	err = db2.Table("guild").Find(&guilds2).Error
	if err != nil {
		return err
	}

	guildMembers1 := make([]*GuildMember, 0)
	guildMembers2 := make([]*GuildMember, 0)

	// load guild member
	err = db1.Table("guild_member").Find(&guildMembers1).Error
	if err != nil {
		return err
	}

	err = db2.Table("guild_member").Find(&guildMembers2).Error
	if err != nil {
		return err
	}

	lastID := uint64(100000)
	length := len(guilds1)
	if length > 0 {
		lastID = guilds1[length-1].ID
	}

	logm.DebugfE("guild lastID [%d]", lastID)

	length2 := len(guilds2)
	//
	guilds1 = append(guilds1, guilds2...)

	guilds2 = make([]*Guild, 0)
	// key preID value afterID
	mapOldNewGuildID := make(map[uint64]uint64)
	// 遍历guilds1从最初的最后一个值+1开始
	if length2 > 0 {
		for index, _ := range guilds1[length:] {
			lastID++

			preID := guilds1[length:][index].ID
			guilds1[length:][index].ID = lastID

			mapOldNewGuildID[preID] = lastID
		}
	}

	// 处理角色
	// 相当于2号数据库里面有角色数据
	if lengthRole < len(roles1) {
		for key, _ := range roles1[lengthRole:] {
			//
			oldID := uint64(roles1[lengthRole:][key].GuildID)
			if oldID == 4294967295 {
				continue
			}

			if newID, ok := mapOldNewGuildID[oldID]; ok {
				roles1[lengthRole:][key].GuildID = uint32(newID)
			} else {
				// 角色身上的工会没有在工会表里面找到怎么办
				roles1[lengthRole:][key].GuildID = 4294967295
				logm.DebugfE("role guild id[%d] not found set guild id[4294967295]", oldID)
			}
		}
	}

	deleteIndex := make([]int, 0)
	for key, value := range guildMembers2 {
		if newID, ok := mapOldNewGuildID[uint64(value.GuildID)]; !ok {
			// 把这条记录删除
			deleteIndex = append(deleteIndex, key)
			logm.DebugfE("guild member guild id[%d] lost", value.GuildID)
		} else {
			guildMembers2[key].GuildID = uint32(newID)
		}
	}

	var count = 0
	for _, value := range deleteIndex {
		guildMembers2 = append(guildMembers2[0:value-count], guildMembers2[value+1-count:]...)
		count++
	}

	t1 := time.Now().Unix()
	////////////////////////////////////////////////////////////////
	// 需要解散的工会逻辑处理
	needDismiss := make(map[uint64]uint64)
	deleteIndex = make([]int, 0)
	// 遍历找出所有需要delete掉的工会
	for i, value := range guilds1 {
		if _, ok := mapRoleIds[value.LeaderID]; ok {
			// 此工会要解散
			needDismiss[value.ID] = value.ID

			deleteIndex = append(deleteIndex, i)
			logm.WarnfE("10此工会因为合并把角色隐藏导致身上的工会需要解散掉 ID[%d]", value.ID)
		}
	}

	// 删除没必要的工会
	count = 0
	for _, value := range deleteIndex {
		guilds1 = append(guilds1[0:value-count], guilds1[value+1-count:]...)
		count++
	}

	// 遍历工会成员表的数据
	deleteIndex = make([]int, 0)
	for key, value := range guildMembers1 {
		if _, ok := mapRoleIds[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, key)

			logm.WarnfE("20此角色因为合区角色被隐藏了,需要从工会成员里面删除掉 roleID[%d] guildID[%d]", value.RoleID, value.GuildID)
			continue
		}

		// 如果上面不存在 在看看在不在解散的工会里面，在里面也删掉
		if _, ok := needDismiss[uint64(value.GuildID)]; ok {
			deleteIndex = append(deleteIndex, key)

			logm.WarnfE("30此角色因为合区导致会长的角色隐藏工会解散掉了,需要从工会成员里面删除掉 roleID[%d] guildID[%d]", value.RoleID, value.GuildID)
		}
	}
	count = 0
	for _, value := range deleteIndex {
		guildMembers1 = append(guildMembers1[0:value-count], guildMembers1[value+1-count:]...)
		count++
	}

	deleteIndex = make([]int, 0)
	for key, value := range guildMembers2 {
		if _, ok := mapRoleIds[value.RoleID]; ok {
			deleteIndex = append(deleteIndex, key)

			logm.WarnfE("40此角色因为合区角色被隐藏了,需要从工会成员里面删除掉 roleID[%d] guildID[%d]", value.RoleID, value.GuildID)
			continue
		}

		// 如果上面不存在 在看看在不在解散的工会里面，在里面也删掉
		if _, ok := needDismiss[uint64(value.GuildID)]; ok {
			deleteIndex = append(deleteIndex, key)

			logm.WarnfE("50此角色因为合区导致会长的角色隐藏工会解散掉了,需要从工会成员里面删除掉 roleID[%d] guildID[%d]", value.RoleID, value.GuildID)
		}
	}

	count = 0
	for _, value := range deleteIndex {
		guildMembers2 = append(guildMembers2[0:value-count], guildMembers2[value+1-count:]...)
		count++
	}

	// 连接起来
	guildMembers1 = append(guildMembers1, guildMembers2...)

	guildMembers2 = make([]*GuildMember, 0)

	// 遍历角色清理身上的工会信息
	for key, _ := range roles1 {
		if roles1[key].GuildID == 4294967295 {
			continue
		}

		if _, ok := mapRoleIds[roles1[key].RoleID]; ok {
			logm.WarnfE("60因为合区导致角色被隐藏了，需要清理身上的工会信息 roleID[%d] guildID[%d]", roles1[key].RoleID, roles1[key].GuildID)
			roles1[key].GuildID = 4294967295
			continue
		}

		if _, ok := needDismiss[uint64(roles1[key].GuildID)]; ok {
			logm.WarnfE("70因为合区导致工会会长被隐藏了，工会解散了，需要清理身上的工会信息 roleID[%d] guildID[%d]", roles1[key].RoleID, roles1[key].GuildID)
			roles1[key].GuildID = 4294967295
		}
	}

	logm.DebugfE("处理清理工会信息耗时 耗时[%d]", time.Now().Unix()-t1)

	// 保存开始
	err = BatchGuildSave(db1, guilds1)
	if err != nil {
		return err
	}

	err = BatchGuildMemberSave(db1, guildMembers1)
	if err != nil {
		return err
	}

	////////////////////////////////////////////////////////////////

	/*
		// batch save guild
		if length2 > 0 {
			err = BatchSave(db1, GuildC, guilds1[length:])
			if err != nil {
				return err
			}
		}

		// batch save guild member
		err = BatchSave(db1, GuildMemberC, guildMembers2)
		if err != nil {
			return err
		}
	*/
	// nil
	// 预定义的标识符
	// 代表一些类型 变量的零值
	return nil
}

func BatchGuildSave(db1 *gorm.DB, guilds []*Guild) error {
	// 先清理1中的数据
	err := db1.Exec("truncate table guild;").Error
	if err != nil {
		return err
	}

	err = BatchSave(db1, GuildC, guilds)
	if err != nil {
		return err
	}

	return nil
}

func BatchGuildMemberSave(db1 *gorm.DB, guildsMem []*GuildMember) error {
	// 先清理1中的member
	err := db1.Exec("truncate table guild_member;").Error
	if err != nil {
		return err
	}

	err = BatchSave(db1, GuildMemberC, guildsMem)
	if err != nil {
		return err
	}

	return nil
}
