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
func HandleGuildRelation(db1, db2 *gorm.DB, roles1 []*RoleData, lengthRole /*这个是roles1合并数据之前的长度*/ int) error {
	guilds1 := make([]*Guild, 0)
	guilds2 := make([]*Guild, 0)

	// load guild
	err := db1.Table("guild").Find(&guilds1).Error
	if err != nil {
		return err
	}

	// load guild
	err = db2.Table("guild").Find(&guilds2).Error
	if err != nil {
		return err
	}

	guildMembers2 := make([]*GuildMember, 0)
	// load guild member

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

	//
	guilds1 = append(guilds1, guilds2...)

	// key preID value afterID
	mapOldNewGuildID := make(map[uint64]uint64)
	// 遍历guilds1从最初的最后一个值+1开始
	if len(guilds2) > 0 {
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
	}

	// batch save guild
	if len(guilds2) > 0 {
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

	// nil
	// 预定义的标识符
	// 代表一些类型 变量的零值
	return nil
}
