package model

import (
	"github.com/gin-gonic/gin"
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/config"
	"gorm.io/gorm"
	"net/http"
	"runtime"
	"time"
)

var CurrentStage = ""
var Count = 0
var Total = 23
var IsStart = false
var Err = ""
var CurrentTime = time.Now().Unix()
var StartTime = time.Now().Unix()

var LastMergerTime = time.Now().Unix()
var LastMergerIp = ""

func ResetDB() {

}

func Reset() {
	CurrentStage = ""
	Count = 0
	IsStart = false
	Err = ""
}

func StartRemoteDB() (error, *gorm.DB, *gorm.DB) {
	mConfig := config.Config.Mysql
	var err error
	var gdb2, gdb4 *gorm.DB

	gdb2, err = NewDB(mConfig.DBUser, mConfig.DBPasswd, mConfig.DBHost2, mConfig.DBNameB, mConfig.DBTablePrefix)
	if err != nil {
		logm.ErrorfE("GDB2连接数据库失败[%s]", err.Error())
		return err, nil, nil
	}

	logm.DebugfE("GDB2连接数据库成功")

	gdb4, err = NewDB(mConfig.DBUser, mConfig.DBPasswd, mConfig.DBHost2, mConfig.DBNameD, mConfig.DBTablePrefix)
	if err != nil {
		logm.ErrorfE("GDB4连接数据库失败[%s]", err.Error())
		return err, nil, nil
	}

	logm.DebugfE("GDB4连接数据库成功")

	SetEnv(true, gdb2)
	return nil, gdb2, gdb4
}

func StartDB(reconnect bool) error {
	var err error
	mConfig := config.Config.Mysql

	if GDB1 == nil {
		GDB1, err = NewDB(mConfig.DBUser, mConfig.DBPasswd, mConfig.DBHost1, mConfig.DBNameA, mConfig.DBTablePrefix)
		if err != nil {
			if reconnect {
				logm.ErrorfE("GDB1重复连接数据库失败[%s]", err.Error())
			} else {
				logm.ErrorfE("GDB1连接数据库失败[%s]", err.Error())
			}
			return err
		}

		if reconnect {
			logm.DebugfE("GDB1重新连接数据库成功")
		} else {
			logm.DebugfE("GDB1连接数据库成功")
		}
	}

	if GDB3 == nil {
		GDB3, err = NewDB(mConfig.DBUser, mConfig.DBPasswd, mConfig.DBHost1, mConfig.DBNameC, mConfig.DBTablePrefix)
		if err != nil {
			if reconnect {
				logm.ErrorfE("GDB3重复连接数据库失败[%s]", err.Error())
			} else {
				logm.ErrorfE("GDB3连接数据库失败[%s]", err.Error())
			}
			return err
		}

		if reconnect {
			logm.DebugfE("GDB3重新连接数据库成功")
		} else {
			logm.DebugfE("GDB3连接数据库成功")
		}
	}

	SetEnv(false, nil)
	return err
}

func SetCurrent(s, err string, delta int64) {
	Err = err
	CurrentTime = time.Now().Unix()

	mConfig := config.Config.Mysql
	// 插入merger_log表 失败日志
	if err != "" {
		InsertMergerLog(CurrentStage, err+"["+mConfig.DBHost2+"]", "合并出错", int(delta))
	} else {
		if CurrentStage != "" {
			InsertMergerLog(CurrentStage, err+"["+mConfig.DBHost2+"]", "合并成功", int(delta))
		}

		if s == Stage_Success {
			InsertMergerLog(s, err+"["+mConfig.DBHost2+"]", "合并成功", int(CurrentTime-StartTime))
		}

		CurrentStage = s
		Count++
	}
}

const (
	Stage_10      = "Start"
	Stage_20      = "Account_common"
	Stage_30      = "Mail"
	Stage_31      = "SelectItem"
	Stage_40      = "item"
	Stage_50      = "Equip"
	Stage_60      = "JiMai"
	Stage_70      = "mail-item"
	Stage_71      = "change_name"
	Stage_80      = "Guild"
	Stage_90      = "Role_data"
	Stage_100     = "BlackList"
	Stage_110     = "Buff"
	Stage_120     = "Enemy"
	Stage_130     = "Friend"
	Stage_140     = "LeftMsg"
	Stage_150     = "MapLimit"
	Stage_160     = "Reward"
	Stage_170     = "RoleExtend"
	Stage_180     = "Skill"
	Stage_190     = "Title"
	Stage_200     = "Sign"
	Stage_210     = "System"
	Stage_Error   = "Error"
	Stage_Success = "success"
)

func Handle() error {
	err := StartDB(false)
	if err != nil {
		logm.ErrorfE("初始化本地数据库失败 err[%s]", err.Error())
		return err
	}

	logm.DebugfE("初始化本地数据库成功")

	err, gdb2, gdb4 := StartRemoteDB()
	if err != nil {
		logm.ErrorfE("初始化远程数据库失败 err[%s]", err.Error())
		return err
	}

	logm.DebugfE("初始化远程数据库连接成功")

	t1 := time.Now().Unix()
	t2 := time.Now().Unix()
	StartTime = time.Now().Unix()

	logm.InfofE("数据合并开始...")

	SetCurrent(Stage_10, "", 0)
	// 远程的数据库数据，合并过来的数据库
	var MinRoleCID uint32
	//var MaxRoleCID uint32
	var MaxRoleID uint32
	var increaseNum uint32

	type MaxStruct struct {
		Max uint32 `json:"max"`
		Min uint32 `json:"min"`
	}
	var result MaxStruct

	err = GDB1.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	MaxRoleID = result.Max

	err = gdb2.Raw("select max(role_id) as max from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	//MaxRoleCID = result.Max

	err = gdb2.Raw("select min(role_id) as min from role_data;").Scan(&result).Error
	if err != nil {
		return err
	}
	MinRoleCID = result.Min

	/*if MaxRoleID >= MaxRoleCID {
		increaseNum = MaxRoleID + 5
	} else {
		increaseNum = MaxRoleCID + 5
	}*/

	if MaxRoleID == 0 {
		increaseNum = 5
	} else {
		increaseNum = MaxRoleID + 5
	}

	MinRoleCID += increaseNum

	err = HandleRelation(GDB1, gdb2, GDB3, gdb4, increaseNum)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta := t2 - CurrentTime
	logm.DebugfE("处理角色ID相关结束...time[%d]", delta)

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
	err = HandleAccountCommon(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理账号表account_common结束 ... time[%d]", delta)

	SetCurrent(Stage_30, "", delta)
	// handle mail 先处理掉
	//err, num := HandleMail(GDB1, GDB2)
	maps, err := HandleMailByIncreaseID(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime
	dd := time.Now().Unix()
	logm.DebugfE("处理邮件mail结束... time[%d]", delta)

	// handle item relative
	err = HandleItemRelation(GDB1, gdb2, GDB3, increaseNum, maps)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - dd
	dd = time.Now().Unix()
	logm.DebugfE("处理物品相关的结束... time[%d]", delta)

	mapRoleIDsRf := make(map[uint32]uint32)
	err = HandleRoleData(GDB1, gdb2, GDB3, gdb4, mapRoleIDsRf)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - dd

	logm.DebugfE("处理角色表相关结束... time[%d]", delta)

	SetCurrent(Stage_100, "", delta)
	err = HandleBlackList(GDB1, gdb2, mapRoleIDsRf)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理黑名单表结束... time[%d]", delta)

	SetCurrent(Stage_110, "", delta)
	err = HandleBuff(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理buff表结束... time[%d]", delta)

	SetCurrent(Stage_120, "", delta)
	err = HandleEnemy(GDB1, gdb2, mapRoleIDsRf)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理敌人表结束... time[%d]", delta)

	SetCurrent(Stage_130, "", delta)
	err = HandleFriend(GDB1, gdb2, mapRoleIDsRf)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理好友表结束... time[%d]", delta)

	SetCurrent(Stage_140, "", delta)
	err = HandleLeftMsg(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理剩余消息表结束 ... time[%d]", delta)

	SetCurrent(Stage_150, "", delta)
	err = HandleMapLimit(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理地图限制表结束... time[%d]", delta)

	SetCurrent(Stage_160, "", delta)
	err = HandleReward(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理奖励表结束... time[%d]", delta)

	SetCurrent(Stage_170, "", delta)
	err = HandleRoleExtend(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理角色额外数据表结束... time[%d]", delta)

	SetCurrent(Stage_180, "", delta)
	err = HandleSkill(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理技能表结束... time[%d]", delta)

	SetCurrent(Stage_190, "", delta)
	err = HandleTitle(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理称号表结束... time[%d]", delta)

	SetCurrent(Stage_200, "", delta)
	err = HandleSign(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理签名表结束... time[%d]", delta)

	SetCurrent(Stage_210, "", delta)
	err = HandleSystem(GDB1, gdb2)
	if err != nil {
		return err
	}
	t2 = time.Now().Unix()
	delta = t2 - CurrentTime

	logm.DebugfE("处理系统变量表结束... time[%d]", delta)

	SetCurrent(Stage_Success, "", delta)

	SetLastMergerInfo(config.Config.Mysql.DBHost2, time.Now().Unix())
	//
	// end ...
	logm.DebugfE("合并成功共结束-消耗时间..... time[%d]", t2-t1)
	return nil
}

type Result struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Total   int    `json:"total"`
	Phase   int    `json:"phase"`
	Err     string `json:"err"`
	IsStart bool   `json:"start"`
}

// http://192.168.16.250:5050/h/merger?ip=127.0.0.1
func HandlerMerger(c *gin.Context) {
	// 输出内存信息
	PrintStatus("合并之前")

	ip := c.Query("ip")

	c.Header("content-type", "text/json")
	isValid := IsValidIp4(ip)
	if !isValid {
		result := Result{
			Code:    999,
			Msg:     "无效的IPv4地址[" + ip + "]",
			Err:     "检测IPv4地址",
			Total:   Total,
			Phase:   Count,
			IsStart: IsStart,
		}

		logm.DebugfE("无效的IPv4地址 ip[%s]", ip)

		c.JSON(http.StatusOK, &result)

		return
	}

	logm.DebugfE("HandlerMerger Start...... ip[%v] path[%s]", ip, c.Request.URL.Path)

	// 已经开始
	if IsStart {
		result := Result{
			Code:    999,
			Msg:     "合并任务正在进行中请稍作等待",
			Err:     "合并已经开始",
			Total:   Total,
			Phase:   Count,
			IsStart: IsStart,
		}

		logm.DebugfE("合并任务正在进行中，重复的请求 ip[%s]", config.Config.Mysql.DBHost2)

		c.JSON(http.StatusOK, &result)

		return
	}

	// 初始化远程IP
	config.Config.Mysql.DBHost2 = ip
	IsStart = true

	result := Result{
		Code:  200,
		Msg:   "合并结束",
		Err:   "",
		Total: Total,
		Phase: Total,
	}

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error: // 运行时错误
				logm.TraceErrorEx("HandlerMerger recover runtime err[%v]", err)
			default: // 非运行时错误
				logm.TraceErrorEx("HandlerMerger recover not runtime err[%v]", err)
			}

			result.Code = 999
			result.Msg = "find panic error"
			result.IsStart = IsStart

			c.JSON(http.StatusOK, &result)
		}
	}()

	err := Handle()

	if err != nil {
		// 设置出错截断
		SetCurrent("", err.Error(), 0)

		result.Code = 999
		result.Msg = err.Error()
		result.Err = Err
		result.Total = Total
		result.Phase = Count
		result.IsStart = IsStart

		logm.ErrorfE("合区出错了 err[%s]", err.Error())
	} else {
		//IsStart = false
		Reset()

		result.IsStart = IsStart

		logm.DebugfE("合区成功了恭喜你")
	}

	runtime.GC()

	// 发送合并结果
	c.JSON(http.StatusOK, &result)
}

func HandleGetSchedule(c *gin.Context) {
	c.Header("content-type", "text/json")

	var result Result
	result.Code = 200
	result.Msg = CurrentStage // 当前阶段
	result.Total = Total
	result.Phase = Count
	result.Err = Err
	result.IsStart = IsStart

	// 返回json数据
	c.JSON(http.StatusOK, &result)
}

func HandleReset(c *gin.Context) {
	c.Header("content-type", "text/json")
	var result Result
	result.Code = 200
	result.Msg = "重置成功，请确保重置之前已经联系了技术处理了错误才可以继续合并操作"

	remoteIP, _ := c.RemoteIP()

	if !IsStart {
		logm.DebugfE("合区没开始，无需重置 remoteip[%s]", remoteIP.String())
		result.Msg = "合区没开始，无需重置"
		c.JSON(http.StatusOK, &result)
		return
	}

	if Err == "" {
		logm.DebugfE("没有发生任何合并错误，无需重置 remoteip[%s]", remoteIP.String())
		result.Msg = "没有发生任何合并错误，无需重置"
		c.JSON(http.StatusOK, &result)
		return
	}

	//IsStart = false
	Reset()
	logm.DebugfE("重置成功了恭喜你 remoteIP[%s]", remoteIP.String())

	// 返回json数据
	c.JSON(http.StatusOK, &result)
}
