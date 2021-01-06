package main

import (
	"encoding/json"
	"fmt"
	"github.com/mx5566/mergernew/model"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"time"
)

func main() {
	model.GDB1, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameA, model.DBTablePrefix)
	model.GDB2, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameB, model.DBTablePrefix)
	model.GDB3, _ = model.NewDB(model.DBUser, model.DBPasswd, model.DBHost, model.DBNameC, model.DBTablePrefix)

	model.SetEnv()

	// 远程的数据库数据，合并过来的数据库
	var MinRoleCID uint32
	var MaxRoleCID uint32
	var MaxRoleID uint32
	var increaseNum uint32

	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct

	err := model.GDB1.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MaxRoleID = result.Max

	err = model.GDB2.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MaxRoleCID = result.Max

	err = model.GDB2.Raw("select min(role_id) as min from role_data;").Scan(&result).Error
	if err != nil {
		panic(err)
	}
	MinRoleCID = result.Min

	if MaxRoleID >= MaxRoleCID {
		increaseNum = MaxRoleID + 5
	} else {
		increaseNum = MaxRoleCID + 5
	}
	MinRoleCID += increaseNum

	// 修改所有合并过来表里面有role_id的字段增加increaseNum
	type TablesRoleID struct {
		TableName string `json:"table_name"`
	}

	// 所有有角色ID（role_id）的字段的表
	var tableRoleID []TablesRoleID

	model.GDB3.Table("COLUMNS").Select("table_name as table_name").Where("TABLE_SCHEMA=? and COLUMN_NAME= ?", model.DBNameB, "role_id").Find(&tableRoleID)

	fmt.Println(tableRoleID)

	// 更新所有上面有role_id的表
	/*for _, table := range tableRoleID {
		model.GDB2.Table(table.TableName).Update("role_id", gorm.Expr("role_id + ?", increaseNum))
	}*/

	//把B表里面所有角色ID需要修改的也增加increaseNum
	err = model.GDB2.Exec("ALTER TABLE `item` ADD INDEX temp_item_o_c (`create_mode`) ").Error
	// model.GDB2.Migrator().CreateIndex()

	if err != nil {
		panic(err)
	}

	type ItemIDs struct {
		ID int64 `gorm:"column:serial"`
	}

	var itemIDs []ItemIDs
	err = model.GDB2.Table("item").Select("serial as serial").Find(&itemIDs).Error
	if err != nil {
		panic(err)
	}

	//model.GDB2.Begin() 好像加了事务还慢了
	model.GDB2.Begin()
	err = model.GDB2.Table("item").Where("container_type_id != ?", 11).Updates(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		model.GDB2.Rollback()
		panic(err)
	}
	model.GDB2.Commit()

	// 所有mode为8的更新increaseNum
	model.GDB2.Begin()
	err = model.GDB2.Table("item").Where("create_mode = ?", 8).Updates(map[string]interface{}{"create_id": gorm.Expr("create_id + ?", increaseNum)}).Error
	if err != nil {
		model.GDB2.Rollback()
		panic(err)
	}
	model.GDB2.Commit()

	//移除索引
	err = model.GDB2.Exec("alter table item drop index temp_item_o_c ;").Error
	if err != nil {
		panic(err)
	}

	/*
		UPDATE guild_hebin_copy SET creater_name_id=creater_name_id+increase_num;
		UPDATE guild_hebin_copy SET leader_id=leader_id+increase_num;
		UPDATE item_del_hebin_copy SET owner_id=owner_id+increase_num;
	*/
	err = model.GDB2.Table("guild").Updates(map[string]interface{}{"creater_name_id": gorm.Expr("creater_name_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("guild").Updates(map[string]interface{}{"leader_id": gorm.Expr("leader_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("item_del").Updates(map[string]interface{}{"owner_id": gorm.Expr("owner_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	/*
		update friend_hebin_copy SET friend_id = friend_id + increase_num;
		update blacklist_hebin_copy set black_id = black_id + increase_num;
		update enemy_hebin_copy set enemy_id = enemy_id + increase_num;
	*/

	err = model.GDB2.Table("friend").Updates(map[string]interface{}{"friend_id": gorm.Expr("friend_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("blacklist").Updates(map[string]interface{}{"black_id": gorm.Expr("black_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("enemy").Updates(map[string]interface{}{"enemy_id": gorm.Expr("enemy_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	/*
		update mail_hebin_copy set recv_role_id = recv_role_id + increase_num where recv_role_id != 4294967295;
		update mail_hebin_copy set send_role_id = send_role_id + increase_num where send_role_id != 4294967295;
	*/

	err = model.GDB2.Table("mail").Where("recv_role_id != ?", 4294967295).Updates(map[string]interface{}{"recv_role_id": gorm.Expr("recv_role_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Table("mail").Where("send_role_id != ?", 4294967295).Updates(map[string]interface{}{"send_role_id": gorm.Expr("send_role_id + ?", increaseNum)}).Error
	if err != nil {
		panic(err)
	}

	// 修改role_name_origin字段
	var target_database = model.GDB1.Migrator().CurrentDatabase()
	var target_database_copy = model.GDB2.Migrator().CurrentDatabase()

	var target_table_name = "role_data"
	var target_table_name_copy = "role_data"
	var column_name = "role_name_origin"

	/*
			select count(*) into @cnt1 FROM information_schema.columns WHERE table_schema = target_database AND table_name = target_table_name AND column_name = target_column_name;
		if @cnt1 = 0 then
				set @st1 = CONCAT('ALTER TABLE ', target_table_name, ' ADD COLUMN ', target_column_name, ' VARCHAR(32) default NULL');
				PREPARE STMT1 FROM @st1;
				EXECUTE STMT1;
				DEALLOCATE PREPARE STMT1;
		end if;

		select count(*) into @cnt2 FROM information_schema.columns WHERE table_schema = target_database AND table_name = target_table_name_copy AND column_name = target_column_name;
		if @cnt2 = 0 then
				set @st2 = CONCAT('ALTER TABLE ', target_table_name_copy, ' ADD COLUMN ', target_column_name, ' VARCHAR(32) default NULL');
				PREPARE STMT2 FROM @st2;
				EXECUTE STMT2;
				DEALLOCATE PREPARE STMT2;
		end if;


	*/

	var c int64 = 0

	// 计算个数
	err = model.GDB3.Table("columns").Where("table_schema = ? and table_name = ? and column_name = ?", target_database, target_table_name, column_name).Select("count(*) as count").Count(&c).Error
	if err != nil {
		panic(err)
	}
	if c == 0 {
		err = model.GDB1.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(32) default NULL;", target_table_name, column_name)).Error
		if err != nil {
			panic(err)
		}
	}

	err = model.GDB3.Table("columns").Where("table_schema = ? and table_name = ? and column_name = ?", target_database_copy, target_table_name_copy, column_name).Select("count(*) as count").Count(&c).Error
	if err != nil {
		panic(err)
	}
	if c == 0 {
		err = model.GDB2.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(32) default NULL;", target_table_name_copy, column_name)).Error
		if err != nil {
			panic(err)
		}
	}

	// 更新两个role_data表
	err = model.GDB1.Exec("update role_data set role_name_origin = role_name where ISNULL(role_name_origin);").Error
	if err != nil {
		panic(err)
	}

	err = model.GDB2.Exec("update role_data set role_name_origin = role_name where ISNULL(role_name_origin);").Error
	if err != nil {
		panic(err)
	}

	// 处理account_common表
	var accounts1 []*model.AccountCommon
	err = model.GDB1.Select("account_id, baibao_yuanbao, total_recharge, yuanbao_recharge, data_ex").Find(&accounts1).Error
	//
	if err != nil {
		panic(err)
	}

	var accounts2 []*model.AccountCommon
	err = model.GDB2.Find(&accounts2).Error
	//
	if err != nil {
		panic(err)
	}

	// 把相同account_id的账号的数据充值元宝合并 data_ex 是个json串需要单独处理
	a2 := make(map[uint32]*model.AccountCommon)
	a1 := make(map[uint32]*model.AccountCommon)

	alreadyAccount := make([]*model.AccountCommon, 0)
	notExistAccount := make([]*model.AccountCommon, 0)

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
		v2.BaiBaoYuanBao += data.BaiBaoYuanBao
		v2.TotalRecharge += data.TotalRecharge
		v2.YuanBaoRecharge += data.YuanBaoRecharge

		// 修改data_ex字段
		d1 := v2.DataEx
		d2 := data.DataEx
		de1 := new(model.DataEx)
		de2 := new(model.DataEx)

		if d1 == "" {
			d1 = "{}"
		}

		if d2 == "" {
			d2 = "{}"
		}

		err := json.Unmarshal([]byte(d1), de1)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(d2), de2)
		if err != nil {
			panic(err)
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
	count := math.Ceil(float64(repeatLength) / float64(model.BaseLength))

	fmt.Println(count)

	// 对于存在的账号数据叠加
	err = model.BatchUpdate(model.GDB1, model.Ac, alreadyAccount)
	if err != nil {
		panic(err)
	}

	// 不存在的数据批量插入
	err = model.BatchSave(model.GDB1, model.Ac, notExistAccount)
	if err != nil {
		panic(err)
	}

	// handle mail
	err = model.HandleMail(model.GDB1, model.GDB2)
}
