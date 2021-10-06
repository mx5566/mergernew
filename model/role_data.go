package model

import (
	"fmt"
	"github.com/mx5566/logm"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RoleData struct {
	AccountID                uint32  `gorm:"column:account_id"`
	RoleID                   uint32  `gorm:"primary_key;column:role_id"`
	RoleName                 string  `gorm:"not null;type:varchar(32);column:role_name"`
	RoleNameCrc              uint32  `gorm:"not null;column:role_name_crc"`
	Sex                      int8    `gorm:"not null;column:sex"`
	HairModelID              uint8   `gorm:"not null;column:hair_model_id"`
	HairColorID              uint8   `gorm:"not null;column:hair_color_id"`
	FaceModelID              uint8   `gorm:"not null;column:face_model_id"`
	FaceDetailID             uint8   `gorm:"not null;column:face_detail_id"`
	DressModelID             uint8   `gorm:"not null;column:dress_model_id"`
	Visualized               uint8   `gorm:"not null;column:visualizeid;default:0"`
	AvatarEquip              []byte  `gorm:"column:avatar_equip"`
	DisplaySet               int8    `gorm:"not null;default:0;column:display_set"`
	MapID                    uint32  `gorm:"not null;column:map_id"`
	X                        float32 `gorm:"not null;column:x"`
	Y                        float32 `gorm:"not null;column:y"`
	Z                        float32 `gorm:"not null;column:z"`
	FaceX                    float32 `gorm:"not null;column:face_x"`
	FaceY                    float32 `gorm:"not null;column:face_y"`
	FaceZ                    float32 `gorm:"not null;column:face_z"`
	RebornMapID              uint32  `gorm:"not null;default:4294967295;column:reborn_map_id"`
	Class                    int8    `gorm:"not null;default:1;column:class"`
	ClassEx                  int8    `gorm:"not null;column:classex"`
	Level                    int32   `gorm:"not null;column:level;default:1"`
	ExpCurLevel              int64   `gorm:"not null;column:exp_cur_level"`
	Hp                       int64   `gorm:"not null;column:hp"`
	Mp                       int64   `gorm:"not null;column:mp"`
	Rage                     int64   `gorm:"not null;column:rage"`
	Endurance                int64   `gorm:"not null;column:endurance;default:-1"`
	Vitality                 int64   `gorm:"not null;column:vitality;default:-1"`
	Injury                   int64   `gorm:"not null;column:injury"`
	Knowledge                int64   `gorm:"not null;column:knowledge"`
	Morale                   int64   `gorm:"not null;column:morale;default:100"`
	Morality                 int64   `gorm:"not null;column:morality"`
	Culture                  uint32  `gorm:"not null;column:culture"`
	Credit                   uint32  `gorm:"not null;column:credit"`
	Identity                 uint8   `gorm:"not null;column:identity"`
	VipPoint                 uint32  `gorm:"not null;column:vip_point"`
	AttAvail                 int64   `gorm:"not null;column:att_avail"`
	TalentAvail              int64   `gorm:"not null;column:talent_avail"`
	PhysiqueAdded            uint16  `gorm:"not null;column:physique_added"`
	StrengthAdded            uint16  `gorm:"not null;column:strength_added"`
	PneumaAdded              uint16  `gorm:"not null;column:pneuma_added"`
	InnerforeAdded           uint16  `gorm:"not null;column:innerforce_added"`
	TechniqueAdded           uint16  `gorm:"not null;column:technique_added"`
	AgilityAdded             uint16  `gorm:"not null;column:agility_added"`
	TalentType1              int8    `gorm:"not null;column:talent_type1;default:-1"`
	TalentType2              int8    `gorm:"not null;column:talent_type2;default:-1"`
	TalentType3              int8    `gorm:"not null;column:talent_type3;default:-1"`
	TalentType4              int8    `gorm:"not null;column:talent_type4;default:-1"`
	TalenrVal1               int16   `gorm:"not null;column:talent_val1"`
	TalenrVal2               int16   `gorm:"not null;column:talent_val2"`
	TalenrVal3               int16   `gorm:"not null;column:talent_val3"`
	TalenrVal4               int16   `gorm:"not null;column:talent_val4"`
	SafeGuardFlag            int8    `gorm:"not null;column:safe_guard_flag;default:1"`
	PkValue                  uint32  `gorm:"not null;column:pk_value"`
	CloseSafeGuardTime       string  `gorm:"not null;column:close_safe_guard_time;char(20)"`
	BagSize                  int16   `gorm:"not null;column:bag_size;default:40"`
	BagGold                  int64   `gorm:"not null;column:bag_gold"`
	BagSilver                int32   `gorm:"not null;column:bag_silver"`
	BagCopper                uint32  `gorm:"not null;column:bag_copper"`
	BagBindGold              int64   `gorm:"not null;column:bag_bind_gold"`
	BagBindSilver            uint32  `gorm:"not null;column:bag_bind_silver"`
	BagBindCopper            uint32  `gorm:"not null;column:bag_bind_copper"`
	BagYuanBao               int64   `gorm:"not null;column:bag_yuanbao"`
	ExchangeVolume           int32   `gorm:"not null;column:exchange_volume"`
	GuildID                  uint32  `gorm:"not null;column:guild_id;default:4294967295"`
	TeamID                   uint32  `gorm:"not null;column:team_id;default:4294967295"`
	TotalTax                 uint32  `gorm:"not null;column:total_tax"`
	RemoteOpenSet            uint32  `gorm:"not null;column:remote_open_set;default:4294967295"`
	CurTitleID               uint16  `gorm:"not null;column:cur_title_id;default:65535"`
	GetMallFreeTime          string  `gorm:"not null;column:get_mall_free_time;char(20)"`
	CreateTime               string  `gorm:"not null;column:create_time;char(20)"`
	LoginTime                string  `gorm:"not null;column:login_time;char(20)"`
	LogoutTime               string  `gorm:"not null;column:logout_time;char(20)"`
	OnlineTime               int32   `gorm:"not null;column:online_time;default:-1"`
	CurOnlineTime            int32   `gorm:"not null;column:cur_online_time"`
	LeaveGuildTime           string  `gorm:"not null;column:leave_guild_time;char(20)"`
	RemoveFlag               int8    `gorm:"not null;column:remove_flag"`
	RemoveTime               uint32  `gorm:"not null;column:remove_time;default:4294967295"`
	ScriptData               []byte  `gorm:"column:script_data"`
	TreasureSum              uint8   `gorm:"not null;column:treasure_sum"`
	StallLevel               uint8   `gorm:"not null;column:stall_level"`
	StallDailyExp            uint32  `gorm:"not null;column:stall_daily_exp"`
	StallCurExp              uint32  `gorm:"not null;column:stall_cur_exp"`
	StallLastTime            uint32  `gorm:"not null;column:stall_last_time;default:4325376"`
	SendMailNum              uint16  `gorm:"not null;column:send_mail_num"`
	MasterID                 uint32  `gorm:"not null;column:master_id"`
	MasterPrenticeForbidTime uint32  `gorm:"not null;column:masterprentice_forbid_time"`
	MapLimitNum              uint32  `gorm:"not null;column:map_limit_num"`
	OwnInstanceID            uint32  `gorm:"not null;column:own_instance_id;default:4294967295"`
	OwnInstanceMapID         uint32  `gorm:"not null;column:own_instance_map_id;default:4294967295"`
	InstanceCreateTime       string  `gorm:"not null;column:instance_create_time;char(20);default:'2000-00-00 00:00:00'"`
	HangNum                  uint8   `gorm:"not null;column:hang_num"`
	IsExp                    uint8   `gorm:"not null;column:is_exp"`
	IsBrotherhood            uint8   `gorm:"not null;column:is_brotherhood"`
	LeaveExp                 int64   `gorm:"not null;column:leave_exp"`
	LeaveBrotherhood         uint32  `gorm:"not null;column:leave_brotherhood"`
	PetPacketNum             uint32  `gorm:"not null;column:pet_packet_num;default:1"`
	RoleHelp                 []byte  `gorm:"column:role_help"`
	RoleTalk                 []byte  `gorm:"column:role_talk"`
	KeyInfo                  []byte  `gorm:"column:key_info"`
	TotalMasterMoral         uint32  `gorm:"not null;column:total_mastermoral"`
	KillNum                  uint32  `gorm:"not null;column:kill_num"`
	GiftGroupID              uint8   `gorm:"not null;column:gift_group_id"`
	GiftStep                 uint8   `gorm:"not null;column:gift_step"`
	GiftID                   uint32  `gorm:"not null;column:gift_id;default:4294967295"`
	GiftLeavingTime          uint32  `gorm:"not null;column:gift_leaving_time"`
	GiftGet                  uint8   `gorm:"not null;column:gift_get"`
	RoleCamp                 int8    `gorm:"not null;column:role_camp"`
	PaiMaiLimit              uint32  `gorm:"not null;column:paimailimit"`
	BankLimit                uint32  `gorm:"not null;column:banklimit"`
	ExBagStep                uint8   `gorm:"not null;column:exbagstep;default:1"`
	ExWareStep               uint8   `gorm:"not null;column:exwarestep;default:1"`
	Vigour                   int64   `gorm:"not null;column:vigour;default:0"`
	TodayOnlineTick          uint32  `gorm:"not null;column:today_online_tick"`
	HistoryVigourCost        int64   `gorm:"not null;column:history_vigour_cost"`
	WareSize                 uint16  `gorm:"not null;column:ware_size;default:30"`
	WareGold                 uint32  `gorm:"not null;column:ware_gold"`
	WareSilver               uint32  `gorm:"not null;column:ware_silver"`
	WareCopper               uint32  `gorm:"not null;column:ware_copper"`
	SignatureName            string  `gorm:"column:signature_name;char(32)"`
	CircleQuest              []byte  `gorm:"column:circle_quest"`
	YuanBaoExchangeNum       uint8   `gorm:"not null;column:yuanbao_exchange_num"`
	AchievementPoint         uint32  `gorm:"not null;column:achievemetn_point"`
	ForbidTalkStart          uint32  `gorm:"not null;column:forbid_talk_start"`
	ForbidTalkEnd            uint32  `gorm:"not null;column:forbid_talk_end"`
	ChangeName               uint8   `gorm:"not null;column:change_name"`
	GraduateNum              int32   `gorm:"not null;column:graduate_num"`
	DestoryEquipCount        uint32  `gorm:"not null;column:destory_equip_count"`
	Cur1v1Score              int32   `gorm:"not null;column:cur_1v1_score"`
	Day1v1Score              int32   `gorm:"not null;column:day_1v1_score"`
	Day1v1Num                int32   `gorm:"not null;column:day_1v1_num"`
	Week1v1Score             int32   `gorm:"not null;column:week_1v1_score"`
	Score1v1Award            int8    `gorm:"not null;column:score_1v1_award"`
	LastChangeNameTime       uint32  `gorm:"not null;column:last_change_name_time;default:4294967295"`
	DeleteRoleGuardTime      uint32  `gorm:"not null;column:delete_role_guard_time;default:4294967295"`
	Exploits                 uint32  `gorm:"not null;column:exploits"`
	CircleQuestRefresh       int32   `gorm:"not null;column:circle_quest_refresh"`
	Exploitslimit            int32   `gorm:"not null;column:exploitslimit"`
	ActiveNum                uint32  `gorm:"not null;column:active_num"`
	ActiveData               []byte  `gorm:"column:active_data"`
	ActiveReceive            []byte  `gorm:"column:active_receive"`
	Justice                  uint32  `gorm:"not null;column:justice;default:0"`
	Purpurdedc               uint32  `gorm:"column:purpuredec;default:0"`
	CircleQuestPerdayNumvber int32   `gorm:"not null;column:circle_quest_perdaynumber;default:20"`
	DayClear                 []byte  `gorm:"column:day_clear"`
	CoolDownReviveCD         int32   `gorm:"not null;column:cooldownrevive_cd;default:30"`
	CircleQuestRefreshDayMax int32   `gorm:"not null;column:circle_quest_refresh_daymax;default:10"`
	ShiHun                   uint32  `gorm:"not null;column:shihun"`
	PerdayHangGetExpTimes    uint32  `gorm:"not null;column:perday_hang_getexp_timems;default:21600000"`
	AchievementNum           uint32  `gorm:"not null;column:achievemetn_num"`
	PetXiuLianSize           uint32  `gorm:"not null;column:pet_xiulian_size;default:1"`
	PerdayVigourGetTotal     uint32  `gorm:"not null;column:perday_vigour_get_total"`
	GuildActiveNum           uint32  `gorm:"not null;column:guild_active_num"`
	GuildActiveData          []byte  `gorm:"column:guild_active_data"`
	GuildActiveReceive       []byte  `gorm:"column:guild_active_receive"`
	GodLevel                 uint32  `gorm:"not null;column:god_level"`
	MasterMoral              int32   `gorm:"not null;column:master_moral"`
	PerdayPrenticeNum        uint32  `gorm:"not null;column:perday_prentice_num"`
	MasterReputionReward     uint8   `gorm:"not null;column:master_repution_reward"`
	CurTitleId2              uint16  `gorm:"not null;column:cur_title_id2;default:65535"`
	CurTitleId3              uint16  `gorm:"not null;column:cur_title_id3;default:65535"`
	InstancePass             uint32  `gorm:"not null;column:instance_pass"`
	ShaodangBeginTime        uint32  `gorm:"not null;column:shaodang_begin_time"`
	ShaodangIndex            int32   `gorm:"not null;column:shaodang_index;default:-1"`
	SpouseID                 uint32  `gorm:"not null;column:spouse_id;default:4294967295"`
	VipLevel                 uint32  `gorm:"not null;column:vip_level"`
	VipDeadLine              string  `gorm:"not null;column:vip_deadline;char(20)"`
	StrdwExternData          []byte  `gorm:"column:strdwExternData"`
	RofbidFlag               int8    `gorm:"not null;column:rofbid_flag"`
	ExternData2              string  `gorm:"not null;column:extern_data2;default:'{}'"`
	RoleNameOrigin           string  `gorm:"column:role_name_origin;char(32)"`
}

type SortRoleData []*RoleData

func (s SortRoleData) Len() int {
	return len(s)
}

func (s SortRoleData) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortRoleData) Less(i, j int) bool {
	if s[i].Level > s[j].Level {
		return true
	} else {
		if s[i].Level == s[j].Level {
			if s[i].TotalTax > s[j].TotalTax {
				return true
			} else if s[i].TotalTax == s[j].TotalTax {
				if s[i].CreateTime < s[j].CreateTime {
					return true
				}
			}
		}
	}

	return false
}

func HandleRoleData(db1, db2, db3, db4 *gorm.DB, mapRoleIDsRf map[uint32]uint32) error {
	// 输出内存信息
	PrintStatus("角色表加载角色数据之前")
	var roles1 []*RoleData
	var roles2 []*RoleData

	t2 := time.Now().Unix()
	delta := t2 - CurrentTime
	SetCurrent(Stage_71, "", delta)

	err := db1.Table("role_data").Select("account_id, role_id, role_name, total_tax, remove_flag, level," +
		"create_time, rofbid_flag, role_name_origin, guild_id").Find(&roles1).Error
	if err != nil {
		return err
	}

	err = db2.Table("role_data").Select("account_id, role_id, role_name, role_name_crc, sex, hair_model_id," +
		"hair_color_id, face_model_id, face_detail_id, dress_model_id, visualizeid, avatar_equip, display_set, map_id, x, y," +
		"z, face_x, face_y, face_z, reborn_map_id, class, classex, level, exp_cur_level, hp, " +
		"mp, rage, endurance, vitality, injury, knowledge, morale, morality, culture, credit, identity," +
		"vip_point, att_avail, talent_avail, physique_added, strength_added, pneuma_added, " +
		"innerforce_added, technique_added, agility_added, talent_type1, talent_type2, talent_type3, " +
		"talent_type4, talent_val1, talent_val2, talent_val3, talent_val4, safe_guard_flag, pk_value," +
		"close_safe_guard_time, bag_size, bag_gold, bag_silver, bag_copper, bag_bind_gold, bag_bind_silver, " +
		"bag_bind_copper, bag_yuanbao, exchange_volume, guild_id, team_id, total_tax, remote_open_set," +
		"cur_title_id, get_mall_free_time, create_time, login_time, logout_time, online_time, cur_online_time," +
		"leave_guild_time, remove_flag, remove_time, treasure_sum, stall_level, stall_daily_exp, stall_cur_exp, " +
		"stall_last_time, send_mail_num, master_id, masterprentice_forbid_time, map_limit_num, own_instance_id," +
		"own_instance_map_id, instance_create_time, hang_num, is_exp, is_brotherhood, leave_exp, leave_brotherhood, " +
		"pet_packet_num, role_talk, total_mastermoral, kill_num, gift_group_id, gift_step, gift_id, gift_leaving_time," +
		"gift_get, role_camp, paimailimit, banklimit, exbagstep, exwarestep, vigour, today_online_tick," +
		"history_vigour_cost, ware_size, ware_gold, ware_silver, ware_copper, signature_name, circle_quest, yuanbao_exchange_num, " +
		"achievemetn_point, forbid_talk_start, forbid_talk_end, change_name, graduate_num, destory_equip_count, cur_1v1_score," +
		"day_1v1_score, day_1v1_num, week_1v1_score, score_1v1_award, last_change_name_time, delete_role_guard_time, " +
		"exploits, circle_quest_refresh, exploitslimit, active_num, active_data, active_receive, justice, purpuredec, circle_quest_perdaynumber," +
		"day_clear, cooldownrevive_cd, circle_quest_refresh_daymax, shihun, perday_hang_getexp_timems, achievemetn_num," +
		"pet_xiulian_size, perday_vigour_get_total, guild_active_num, guild_active_data, guild_active_receive, " +
		"god_level, master_moral, perday_prentice_num, master_repution_reward, cur_title_id2, " +
		"cur_title_id3, instance_pass, shaodang_begin_time, shaodang_index, spouse_id, vip_level, vip_deadline," +
		"rofbid_flag, role_name_origin").Find(&roles2).Error
	if err != nil {
		return err
	}

	// 输出内存信息
	PrintStatus("角色表加载角色数据之后")
	// 记录下来所有RofbidFlag == 2 和 == 1 的角色id

	//length := len(roles1)
	// 内存执行效率还是高
	// 把没有初始化的名字全部初始化
	for key, value := range roles1 {
		//logm.DebugfE("name[%s] size[%d] char[%d]", value.RoleNameOrigin, len(value.RoleNameOrigin), len([]rune(value.RoleNameOrigin)))
		if value.RoleNameOrigin == "" {
			roles1[key].RoleNameOrigin = value.RoleName
		}

		roles1[key].TotalTax = 1

		if roles1[key].RofbidFlag <= 1 {
			roles1[key].RofbidFlag = 0
		} else {
			mapRoleIDsRf[value.RoleID] = value.RoleID
		}
	}

	//roles1Temp := roles1

	// 把没有初始化的名字全部初始化
	for key, value := range roles2 {
		if value.RoleNameOrigin == "" {
			roles2[key].RoleNameOrigin = value.RoleName
		}

		roles2[key].TotalTax = 0

		if roles2[key].RofbidFlag <= 1 {
			roles2[key].RofbidFlag = 0
		}
	}

	// 最初的roles1的角色个数， 用来拆分数组的
	lengthRole1 := len(roles1)
	lengthRole2 := len(roles2)
	// 用户内存合并
	roles1 = append(roles1, roles2...)

	roles2 = make([]*RoleData, 0)

	// 同一个账号下面没有删除合没有禁用的角色对应的roles1里面的索引存储下来
	mapAccountRoles := make(map[uint32][]int)
	// 处理之前的老数据
	for key, value := range roles1 {
		roles1[key].RoleNameOrigin = strings.Replace(value.RoleNameOrigin, "_", "", -1)

		length := len(roles1[key].RoleNameOrigin)
		if length > 21 {
			// 截掉一个字符
			r := []rune(roles1[key].RoleNameOrigin)
			r = r[0 : len(r)-1]

			name := string(r)
			// 转换回来
			roles1[key].RoleNameOrigin = name
			roles1[key].RoleName = name
		} else {
			roles1[key].RoleName = roles1[key].RoleNameOrigin
		}

		// 个数+1
		if value.RemoveFlag == 0 && value.RofbidFlag == 0 {
			mapAccountRoles[value.AccountID] = append(mapAccountRoles[value.AccountID], key)
		}
	}

	//////////////////////////////////////////
	// 同一账号下角色超过三个的角色保留处理
	for _, value := range mapAccountRoles {
		if len(value) <= 3 {
			continue
		}

		rds := make([]*RoleData, 0)
		map1 := make(map[uint32]int)
		for _, v := range value {
			rds = append(rds, roles1[v])
			map1[roles1[v].RoleID] = v
		}

		// 排序 出去前三个之后剩下的全部把他rofbid_flag=1
		sort.Sort(SortRoleData(rds))

		// 把剩余的角色置为禁用
		r := rds[3:]
		for _, v := range r {
			index, _ := map1[v.RoleID]
			roles1[index].RofbidFlag = 1

			// 记录封禁的号
			mapRoleIDsRf[v.RoleID] = v.RoleID
		}
	}
	//////////////////////////////////////////

	err = RenameRole(roles1)
	if err != nil {
		return err
	}

	err = RenameRole(roles1)
	if err != nil {
		return err
	}

	// delete role_name_crc index
	var target_database = db1.Migrator().CurrentDatabase()
	var target_table_name = "role_data"
	var c int64 = 0
	// 计算个数
	err = db3.Table("statistics").Where("index_schema = ? and table_name = ? and index_name = ?", target_database, target_table_name, "role_name_crc").Select("count(*) as count").Count(&c).Error
	if err != nil {
		return err
	}

	if c > 0 {
		// 删除role_name_crc的索引
		str := fmt.Sprintf("drop index role_name_crc on role_data;")
		err = db1.Exec(str).Error
		if err != nil {
			return err
		}
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	SetCurrent(Stage_80, "", delta)

	// handle role guild
	err = HandleGuildRelation(db1, db2, roles1, lengthRole1, mapRoleIDsRf)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	PrintStatus("角色表保存角色数据之前")

	logm.DebugfE("handle guild relation ok ..... time[%d]", delta)

	SetCurrent(Stage_90, "", delta)
	err = SaveRoleData(db1, db2, roles1, lengthRole1, lengthRole2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("save role data ok ..... time[%d]", delta)

	// 输出内存信息
	PrintStatus("角色表保存角色数据之后")
	//
	return nil
}

// function RenameRole handle repeat name rename
// @param []*RoleData
// It Return error if no error is nil else error
func RenameRole(roles1 []*RoleData) error {
	// 执行第一遍
	mapNamesRoles := make(map[string][]*RoleData)
	mapRoleToIndex := make(map[uint32]int)
	// 同意账号下所有的角色
	//mapAllAccountRoles := make(map[uint32][]int)
	for key, value := range roles1 {
		// 个数+1
		mapNamesRoles[roles1[key].RoleName] = append(mapNamesRoles[roles1[key].RoleName], value)
		// 角色id对应的索引
		mapRoleToIndex[roles1[key].RoleID] = key
	}

	// 找到重名个数>1
	// 然后找到这个名字的所有角色进行改名处理
	for _, value := range mapNamesRoles {
		if len(value) <= 1 {
			continue
		}

		// 重名超过一个的角色
		rds := make([]*RoleData, 0)
		rds = append(rds, value...)

		// 排序
		sort.Sort(SortRoleData(rds))

		// 改名逻辑

		for k1, v1 := range rds {

			k := k1 + 1
			// 第一个不处理
			if k == 1 {
				continue
			}

			name := v1.RoleName
			// 追加一个字符 如果字节大于21 截断前六个字符，然后追加a-z A-Z
			// 反之直接追加字符
			// 重名的 从第二个开始一直追加
			if len(v1.RoleName+"A") > 21 {
				// 截掉一个字符
				r := []rune(name)
				r = r[0 : len(r)-1]
				name = string(r)
			}

			var i rune
			if k > 27 {
				if k > 53 {
					i = 48
				} else {
					i = rune(k + 69)
				}
			} else {
				i = rune(k + 63)
			}

			// 转换回来
			roles1[mapRoleToIndex[v1.RoleID]].RoleName = name + string(i)
		}
	}
	return nil
}

func SaveRoleData(db1, db2 *gorm.DB, roles1 []*RoleData, length1, length2 int) error {
	PrintStatus("批量更新db1中角色数据之前")
	t1 := time.Now().Unix()
	// 1、先更新db1中role_data的数据
	err := BatchUpdate(db1, RoleDataUp1C, roles1[0:length1])
	if err != nil {
		return err
	}
	delta := time.Now().Unix() - t1
	t1 = time.Now().Unix()
	logm.DebugfE("批量更新角色表数据 耗时[%d]", delta)
	PrintStatus("批量更新db1中角色数据之后")

	// 保存之前线截断role_data表
	//	//err := db1.Table("role_data").Exec("truncate table role_data;").Error
	//	//if err != nil {
	//	//	return err
	//	//}

	// 2、批量保存db2中的role_data数据到db1中
	// 把2中的不包含大二进制数据的字段先插入数据库
	PrintStatus("批量保存db2中角色数据到db1之前")
	err = BatchSave(db1, RoleDataC, roles1[length1:])
	if err != nil {
		return err
	}
	delta = time.Now().Unix() - t1
	t1 = time.Now().Unix()
	logm.DebugfE("批量保存db2中的角色数据 耗时[%d]", delta)
	PrintStatus("批量保存db2中角色数据到db1之后")

	PrintStatus("批量更新db2中角色数据到db1之前")
	// 3、分批把2中的大的二进制数据加载进来更新到1中
	err = HandleUpdateRoleDataBlob(db1, db2, length2)
	if err != nil {
		return err
	}
	delta = time.Now().Unix() - t1
	logm.DebugfE("批量更新保存db2大的二进制数据 耗时[%d]", delta)
	PrintStatus("批量更新db2中角色数据到db1之后")

	return nil
}

func HandleUpdateRoleDataBlob(db1, db2 *gorm.DB, length int) error {
	base := 10000

	tableCount := math.Ceil(float64(length) / float64(base))

	for i := 0; i < int(tableCount); i++ {
		// 输出内存信息
		PrintStatus("角色表分批保存额外大二进制数据之前" + strconv.Itoa(i))

		var rolesT []*RoleData
		err := db2.Table("role_data").Select("role_id, script_data, role_help, " +
			"key_info, strdwExternData, extern_data2").Offset(i * base).Limit(base).Find(&rolesT).Error
		if err != nil {
			return err
		}

		err = BatchUpdate(db1, RoleDataUp2C, rolesT)
		if err != nil {
			logm.ErrorfE("批量更新角色表失败[%s]", err.Error())
			return err
		}

		rolesT = make([]*RoleData, 0)
		// 输出内存信息
		PrintStatus("角色表分批保存额外大二进制数据之后" + strconv.Itoa(i))
	}

	return nil
}
