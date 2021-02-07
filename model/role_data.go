package model

import (
	"fmt"
	"github.com/mx5566/logm"
	"gorm.io/gorm"
	"sort"
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
	ExpCurLevel              uint32  `gorm:"not null;column:exp_cur_level"`
	Hp                       int64   `gorm:"not null;column:hp"`
	Mp                       int64   `gorm:"not null;column:mp"`
	Rage                     uint32  `gorm:"not null;column:rage"`
	Endurance                int32   `gorm:"not null;column:endurance;default:-1"`
	Vitality                 int32   `gorm:"not null;column:vitality;default:-1"`
	Injury                   uint32  `gorm:"not null;column:injury"`
	Knowledge                uint32  `gorm:"not null;column:knowledge"`
	Morale                   uint32  `gorm:"not null;column:morale;default:100"`
	Morality                 uint32  `gorm:"not null;column:morality"`
	Culture                  uint32  `gorm:"not null;column:culture"`
	Credit                   uint32  `gorm:"not null;column:credit"`
	Identity                 uint8   `gorm:"not null;column:identity"`
	VipPoint                 uint32  `gorm:"not null;column:vip_point"`
	AttAvail                 uint16  `gorm:"not null;column:att_avail"`
	TalentAvail              uint32  `gorm:"not null;column:talent_avail"`
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
	BagGold                  int32   `gorm:"not null;column:bag_gold"`
	BagSilver                int32   `gorm:"not null;column:bag_silver"`
	BagCopper                uint32  `gorm:"not null;column:bag_copper"`
	BagBindGold              uint32  `gorm:"not null;column:bag_bind_gold"`
	BagBindSilver            uint32  `gorm:"not null;column:bag_bind_silver"`
	BagBindCopper            uint32  `gorm:"not null;column:bag_bind_copper"`
	BagYuanBao               int32   `gorm:"not null;column:bag_yuanbao"`
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
	LeaveExp                 uint32  `gorm:"not null;column:leave_exp"`
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
	Vigour                   uint16  `gorm:"not null;column:vigour"`
	TodayOnlineTick          uint32  `gorm:"not null;column:today_online_tick"`
	HistoryVigourCost        uint32  `gorm:"not null;column:history_vigour_cost"`
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
	ExternData2              string  `gorm:"not null;column:extern_data2"`
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
	return s[i].Level > s[j].Level && s[i].TotalTax > s[j].TotalTax && s[i].CreateTime < s[j].CreateTime
}

func HandleRoleData(db1, db2, db3, db4 *gorm.DB) error {
	var roles1 []*RoleData
	var roles2 []*RoleData

	err := db1.Find(&roles1).Error
	if err != nil {
		return err
	}

	err = db2.Find(&roles2).Error
	if err != nil {
		return err
	}

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

	// 用户内存合并
	roles1 = append(roles1, roles2...)

	// 同一个账号下面没有删除合没有禁用的角色对应的roles1里面的索引存储下来
	mapAccountRoles := make(map[uint32][]int)
	// 处理之前的老数据
	for key, value := range roles1 {
		roles1[key].RoleNameOrigin = strings.Replace(value.RoleNameOrigin, "_", "", -1)

		length := len(value.RoleNameOrigin)
		if length > 21 {
			// 截掉一个字符
			r := []rune(value.RoleNameOrigin)
			r = r[0 : len(r)-1]

			name := string(r)
			// 转换回来
			roles1[key].RoleNameOrigin = name
			roles1[key].RoleName = name
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
		str := fmt.Sprintf("drop index role_name_crc on role_data")
		err = db1.Exec(str).Error
		if err != nil {
			return err
		}
	}
	t2 := time.Now().Unix()
	delta := t2 - CurrentTime
	SetCurrent(Stage_80, "", delta)

	// handle role guild
	err = HandleGuildRelation(db1, db2, roles1, lengthRole1)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle guild relation ok ..... time[%d]", delta)

	SetCurrent(Stage_90, "", delta)
	err = SaveRoleData(db1, roles1)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("save role data ok ..... time[%d]", delta)

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

func SaveRoleData(db1 *gorm.DB, roles1 []*RoleData) error {
	// 保存之前线截断role_data表
	err := db1.Table("role_data").Exec("truncate table role_data;").Error
	if err != nil {
		return err
	}

	// 先清理db1的role_data的所有数据
	err = BatchSave(db1, RoleDataC, roles1)
	if err != nil {
		return err
	}

	return nil
}
