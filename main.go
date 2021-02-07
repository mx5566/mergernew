package main

import (
	"github.com/mx5566/logm"
	"github.com/mx5566/mergernew/merger"
	"github.com/mx5566/mergernew/model"
	_ "gorm.io/driver/mysql"
	_ "net/http/pprof"
)

// 性能查看
// https://blog.csdn.net/wangdenghui2005/article/details/99119941
func main() {
	// init log module
	logm.Init("mergernew", map[string]string{"errFile": "game_err.log", "logFile": "game.log"}, "debug")

	err := model.StartDB()
	if err != nil {
		logm.PanicfE(err.Error())
	}

	logm.InfofE("Start DB ok...")

	merger.Run()
}
