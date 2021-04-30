package route

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mx5566/mergernew/config"
	"github.com/mx5566/mergernew/model"
)

var R = Default()

var GinR = gin.Default()

//http://liumurong.org/2019/12/gin_pprof/
// register route
func Init() {
	if config.Config.Mode == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	pprof.Register(GinR) // 性能
	//R.GET("/debug/pprof/", pprof.Index)
	//R.GET("/debug/pprof/cmdline", pprof.Cmdline)
	//R.GET("/debug/pprof/profile", pprof.Profile)
	//R.GET("/debug/pprof/symbol", pprof.Symbol)
	//R.GET("/debug/pprof/trace", pprof.Trace)

	//R.GET("/h/merger", model.HandlerMerger)
	//R.GET("/h/schedule", model.HandleGetSchedule)

	// gin 路由注册
	GinR.GET("/h/merger", model.HandlerMerger)
	GinR.GET("/h/schedule", model.HandleGetSchedule)
	GinR.GET("/h/reset", model.HandleReset)

}

// route
// <scheme>://<user>:<password>@<host>:<port>/<path>;<params>?<query>#<fragment>
// https://www.jd.com/shoes?cu=true&utm_source=baidu
// http://wwww.aa.com:8080/A/B?a=1&b=2
// /A/B 	-> path
// host 	-> wwww.aa.com
// port 	-> 8080
// query	-> a=1&b=2
// http -> method -> uri(完整的地址) -> path匹配 -> handler
