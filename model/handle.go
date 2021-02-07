package model

import (
	"encoding/json"
	"fmt"
	"github.com/mx5566/logm"
	"net/http"
	"runtime"
	"time"
)

func StartDB() error {
	GDB1, _ = NewDB(DBUser, DBPasswd, DBHost1, DBNameA, DBTablePrefix)
	GDB2, _ = NewDB(DBUser, DBPasswd, DBHost2, DBNameB, DBTablePrefix)
	GDB3, _ = NewDB(DBUser, DBPasswd, DBHost1, DBNameC, DBTablePrefix)
	GDB4, _ = NewDB(DBUser, DBPasswd, DBHost2, DBNameD, DBTablePrefix)

	SetEnv()
	return nil
}

var CurrentStage = ""
var Count = 0
var Total = 22
var IsStart = false
var Err = ""
var CurrentTime = time.Now().Unix()

func SetCurrent(s, err string, delta int64) {
	Err = err
	CurrentTime = time.Now().Unix()

	// 插入merger_log表 失败日志
	if err != "" {
		InsertMergerLog(CurrentStage, err+"["+DBHost2+"]", "合并出错", int(delta))
	} else {
		if CurrentStage != "" {
			InsertMergerLog(CurrentStage, err+"["+DBHost2+"]", "合并成功", int(delta))
		}
		CurrentStage = s
		Count++
	}
}

const (
	Stage_10    = "Start"
	Stage_20    = "Account_common"
	Stage_30    = "Mail"
	Stage_31    = "SelectItem"
	Stage_40    = "Equip"
	Stage_50    = "JiMai"
	Stage_60    = "mail-item"
	Stage_70    = "item"
	Stage_80    = "Guild"
	Stage_90    = "Role_data"
	Stage_100   = "BlackList"
	Stage_110   = "Buff"
	Stage_120   = "Enemy"
	Stage_130   = "Friend"
	Stage_140   = "LeftMsg"
	Stage_150   = "MapLimit"
	Stage_160   = "Reward"
	Stage_170   = "RoleExtend"
	Stage_180   = "Skill"
	Stage_190   = "Title"
	Stage_200   = "Sign"
	Stage_210   = "System"
	Stage_Error = "Error"
)

func Handle() error {
	t1 := time.Now().Unix()
	t2 := time.Now().Unix()

	logm.InfofE("mergernew start...")

	SetCurrent(Stage_10, "", 0)
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

	err := GDB1.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	MaxRoleID = result.Max

	err = GDB2.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	MaxRoleCID = result.Max

	err = GDB2.Raw("select min(role_id) as min from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	MinRoleCID = result.Min

	if MaxRoleID >= MaxRoleCID {
		increaseNum = MaxRoleID + 5
	} else {
		increaseNum = MaxRoleCID + 5
	}
	MinRoleCID += increaseNum

	err = HandleRelation(GDB1, GDB2, GDB3, GDB4, increaseNum)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta := t2 - CurrentTime
	logm.DebugfE("handle cluster ok ... time[%d]", delta)

	// 用来测试时间
	// 不利用主键去count 利用第二索引去count
	/*var c1 int64 = 0
	err = model.GDB2.Table("item").Select("count(container_type_id) as count").Where("container_type_id > ?", 0).Count(&c1).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("最大个数 ", c1)*/

	SetCurrent(Stage_20, "", delta)
	// handle account_common
	err = HandleAccountCommon(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle account_common ok ... time[%d]", delta)

	SetCurrent(Stage_30, "", delta)
	// handle mail 先处理掉
	//err, num := HandleMail(GDB1, GDB2)
	maps, err := HandleMailByIncreaseID(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle mail ok ... time[%d]", delta)

	il := time.Now().Unix()
	// handle item relative
	err = HandleItemRelation(GDB1, GDB2, increaseNum, maps)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - il
	logm.DebugfE("handle item relation ok ... time[%d]", delta)

	rd := time.Now().Unix()
	err = HandleRoleData(GDB1, GDB2, GDB3, GDB4)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - rd
	logm.DebugfE("handle role data ok ..... time[%d]", delta)

	delta = t2 - CurrentTime
	SetCurrent(Stage_100, "", delta)
	err = HandleBlackList(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle blacklist ok ... time[%d]", delta)

	SetCurrent(Stage_110, "", delta)
	err = HandleBuff(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle buff ok ... time[%d]", delta)

	SetCurrent(Stage_120, "", delta)
	err = HandleEnemy(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle enemy ok ... time[%d]", delta)

	SetCurrent(Stage_130, "", delta)
	err = HandleFriend(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle friend ok ... time[%d]", delta)

	SetCurrent(Stage_140, "", delta)
	err = HandleLeftMsg(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle leftmsg ok ... time[%d]", delta)

	SetCurrent(Stage_150, "", delta)
	err = HandleMapLimit(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle maplimit ok ... time[%d]", delta)

	SetCurrent(Stage_160, "", delta)
	err = HandleReward(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle reward ok ... time[%d]", delta)

	SetCurrent(Stage_170, "", delta)
	err = HandleRoleExtend(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle role_extend ok ... time[%d]", delta)

	SetCurrent(Stage_180, "", delta)
	err = HandleSkill(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle skill ok ... time[%d]", time.Now().Unix()-t2)

	SetCurrent(Stage_190, "", delta)
	err = HandleTitle(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle title ok ... time[%d]", time.Now().Unix()-t2)

	SetCurrent(Stage_200, "", delta)
	err = HandleSign(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle sign ok ... time[%d]", delta)

	SetCurrent(Stage_210, "", delta)
	err = HandleSystem(GDB1, GDB2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	logm.DebugfE("handle system ok ... time[%d]", delta)

	SetCurrent("success", "", t2-t1)
	//
	// end ...
	logm.DebugfE("mergernew success 消耗时间..... time[%d]", t2-t1)
	return nil
}

type Result struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Total int    `json:"total"`
	Phase int    `json:"phase"`
	Err   string `json:"err"`
}

func HandlerMerger(w http.ResponseWriter, req *http.Request) {
	logm.DebugfE("HandlerMerger Start......")

	//InsertMergerLog("1", "testerrorlog", "ok", 1)

	w.Header().Set("content-type", "text/json")

	if IsStart {
		result := Result{
			Code: 999,
			Msg:  "重复的合并操作，无效",
		}

		ret, err := json.Marshal(&result)
		if err != nil {
			logm.ErrorfE("handler recover marshal err[%s] data[%v]", err.Error(), result)
		} else {
			_, err = w.Write(ret)
			if err != nil {
				logm.ErrorfE("handler recover write data err[%s] data[%v]", err.Error(), ret)
			}
		}

		return
	}

	IsStart = true
	err := Handle()
	if err == nil {
		var a = []int{1, 2, 3}
		if a[5] == 1 {

		}
	}

	result := Result{
		Code: 200,
		Msg:  "",
	}

	defer func() {
		if err := recover(); err != nil {

			fmt.Println(err)

			switch err.(type) {
			case runtime.Error: // 运行时错误
				logm.TraceErrorEx("HandlerMerger recover runtime err[%v]", err)
			default: // 非运行时错误
				logm.TraceErrorEx("HandlerMerger recover not runtime err[%v]", err)
			}

			result.Code = 999
			result.Msg = "find panic error"

			ret, err := json.Marshal(&result)
			if err != nil {
				logm.ErrorfE("handler recover marshal err[%s] data[%v]", err.Error(), result)
			} else {
				_, err = w.Write(ret)
				if err != nil {
					logm.ErrorfE("handler recover write data err[%s] data[%v]", err.Error(), ret)
				}
			}

			IsStart = false
		}
	}()

	if err != nil {
		result.Code = 999
		result.Msg = err.Error()

		logm.ErrorfE("handler merger handle err[%s]", err.Error())

		// 设置出错截断
		SetCurrent("", err.Error(), 0)
	}

	ret, err := json.Marshal(&result)
	if err != nil {
		logm.ErrorfE("handler merger marshal err[%s] data[%v]", err.Error(), result)
	} else {
		_, err = w.Write(ret)
		if err != nil {
			logm.ErrorfE("handler merger write data err[%s] data[%v]", err.Error(), ret)
		}
	}

	IsStart = false
}

func HandleGetSchedule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/json")

	type Schedule struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var result Result
	result.Code = 200
	result.Msg = CurrentStage // 当前阶段
	result.Total = Total
	result.Phase = Count
	result.Err = Err
	ret, err := json.Marshal(&result)

	if err != nil {
		logm.ErrorfE("handler get schedule marshal err[%s] data[%v]", err.Error(), result)
	} else {
		_, err = w.Write(ret)
	}
}
