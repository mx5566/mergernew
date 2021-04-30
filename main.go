package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/config"
	"github.com/mx5566/mergernew/merger"
	"github.com/mx5566/mergernew/route"
	_ "gorm.io/driver/mysql"
	"log"
	_ "net/http/pprof"
	"os"
)

var serviceConfig = &service.Config{
	Name:        "合区服务",
	DisplayName: "merger_service",
	Description: "游戏服务器数据合并服务",
}

type FlagParams struct {
	Dir      string // 路径
	Operator string // 操作类型 install start1 stop uninstall restart
	Daemon   bool   // 是否后台执行
}

var Params FlagParams

func init() {
	InitFlag()
}

func InitFlag() {
	flag.StringVar(&Params.Dir, "dir", "", "可执行程序的目录 install的cmd下为必填的选项")
	flag.StringVar(&Params.Operator, "cmd", "", "操作的命令 install uninstall start1 stop restart")
	flag.BoolVar(&Params.Daemon, "d", false, "是否后台执行 true后台 false前台(其他参数无需填写)")
}

// 性能查看
// https://blog.csdn.net/wangdenghui2005/article/details/99119941
func main() {
	// 解析命令行参数
	flag.Parse()

	// mergernew.exe -d true -dir C:\mergernew\ -cmd install
	// mergernew.exe -d true  -cmd start1
	if Params.Daemon {
		if Params.Operator == "" && Params.Dir == "" {
			Params.Daemon = false
		}
	}

	if Params.Operator == "install" {
		if Params.Dir == "" {
			os.Exit(0)
		}

		serviceConfig.Arguments = append(serviceConfig.Arguments, "-dir")
		serviceConfig.Arguments = append(serviceConfig.Arguments, Params.Dir)
	}

	// 切换目录因为如果服务安装之后(管理员权限) 目录就会在system32目录下面
	if Params.Dir != "" {
		err := os.Chdir(Params.Dir)
		if err != nil {
			fmt.Printf("切换目录失败[%s]", err.Error())
			os.Exit(0)
		}
	}

	cur, _ := os.Getwd()
	fmt.Println("当前的目录->", cur)

	// init log module
	logm.Init("mergernew", map[string]string{"errFile": "merger_err.log", "logFile": "merger.log"}, "debug")
	logm.DebugfE("日志模块初始化结束")

	logm.DebugfE("当前目录[%s] 参数[%v]", cur, Params)

	// 构建服务对象
	prog := &Program{}

	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	logm.DebugfE("", os.Args)
	if !Params.Daemon || Params.Operator == "" {
		err = s.Run()
		if err != nil {
			logm.ErrorfE(err.Error())
		}
		return
	}

	cmd := Params.Operator

	fmt.Println("服务命令", cmd)

	if cmd == "install" {
		err = s.Install()
		if err != nil {
			// 文本日志
			logm.ErrorfE("安装服务出错[%s]", err.Error())
			return
		}
		logm.DebugfE("安装成功")
		return
	}

	if cmd == "uninstall" {
		status, err := s.Status()
		if status == service.StatusRunning {
			err = s.Stop()
		}

		err = s.Uninstall()
		if err != nil {
			// 文本日志
			logm.ErrorfE("卸载服务出错[%s]", err.Error())
			return
		}
		logm.DebugfE("卸载成功")
		return
	}

	if cmd == "start1" {
		err = s.Start()
		if err != nil {
			// 文本日志
			logm.ErrorfE("启动服务出错[%s]", err.Error())
			return
		}

		logm.DebugfE("启动成功")
		return
	}

	if cmd == "stop" {
		err = s.Stop()
		if err != nil {
			// 文本日志
			logm.ErrorfE("停止服务出错[%s]", err.Error())
			return
		}

		logm.DebugfE("停止成功")
		return
	}

	if cmd == "restart" {
		status, err := s.Status()
		if status == service.StatusRunning {
			err = s.Restart()
		} else {
			err = s.Start()
		}

		if err != nil {
			// 文本日志
			logm.ErrorfE("重启服务出错[%s]", err.Error())
			return
		}
		logm.DebugfE("重启成功")
		return
	}

	if cmd == "status" {
		status, err := s.Status()
		if err != nil {
			// 文本日志
			logm.ErrorfE("查看服务状态出错[%s]", err.Error())
			return
		}

		logm.DebugfE("服务状态[%d]", status)
		return
	}

}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	logm.DebugfE("开始服务")
	err := Init()
	if err != nil {
		return err
	}

	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	logm.DebugfE("停止服务")

	return nil
}

func (p *Program) run() {
	logm.DebugfE("运行服务run")

	// 启动监听服务
	err := merger.Run()
	if err != nil {
		logm.ErrorfE("启动端口监听失败[%s]", err.Error())
		return
	}
}

func Init() error {
	// 配置文件加载
	err := config.ReadConfig()
	if err != nil {
		logm.ErrorfE("读取配置文件失败[%s]", err.Error())
		return err
	}
	logm.DebugfE("读取配置文件结束")

	// 初始化路由信息
	route.Init()

	logm.DebugfE("路由信息初始话结束")

	//err = model.StartDB(false)
	//if err != nil {
	//	logm.ErrorfE("初始化数据库失败[%s]", err.Error())
	//	logm.DebugfE("定时器开始")
	//
	//	retTime := time.After(time.Second * 2)
	//	<-retTime
	//	logm.DebugfE("定时器结束")
	//	_ = model.StartDB(true)
	//}
	//
	//logm.DebugfE("初始化数据库结束")

	return nil
}
