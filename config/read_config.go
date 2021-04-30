package config

import (
	"github.com/go-ini/ini"
	"os"
	"runtime"
)

var (
	CONFIG_FILE = "config.ini"
)

var Config MergerConfig

type MergerConfig struct {
	Mode    string // 开发模式
	Port    uint16 // 服务监听的端口
	OpenMem bool   // 是否开启内存打印
	Mysql   MysqlConfig
}

type MysqlConfig struct {
	DBType        string
	DBUser        string
	DBPasswd      string
	DBHost1       string
	DBHost2       string
	DBNameA       string
	DBNameB       string
	DBNameC       string
	DBNameD       string
	DBTablePrefix string
}

func ReadConfig() error {
	cur, _ := os.Getwd()

	oo := runtime.GOOS
	if oo == "windows" {
		cur += "\\"
	} else {
		cur += "/"
	}

	cfg, err := ini.Load(cur + CONFIG_FILE)
	if err != nil {
		return err
	}

	Config.Mode = cfg.Section("").Key("app_mode").String()

	var port uint
	port, err = cfg.Section("").Key("port").Uint()
	if err != nil {
		port = 5050
	}

	Config.Port = uint16(port)

	Config.OpenMem, err = cfg.Section("").Key("open_mem").Bool()
	if err != nil {
		Config.OpenMem = false
	}

	err = ReadMysqlConfig(cfg.Section("database"), &Config.Mysql)
	if err != nil {
		return err
	}

	return nil
}

func ReadMysqlConfig(section *ini.Section, mc *MysqlConfig) error {
	mc.DBType = section.Key("DBType").String()
	mc.DBUser = section.Key("DBUser").String()
	mc.DBPasswd = section.Key("DBPasswd").String()
	mc.DBHost1 = section.Key("DBHost1").String()
	mc.DBHost2 = section.Key("DBHost2").String()
	mc.DBNameA = section.Key("DBNameA").String()
	mc.DBNameB = section.Key("DBNameB").String()
	mc.DBNameC = section.Key("DBNameC").String()
	mc.DBNameD = section.Key("DBNameD").String()
	mc.DBTablePrefix = ""

	return nil
}
