package route

import (
	"github.com/mx5566/mergernew/model"
	"net/http/pprof"
)

var R = Default()

// register route
func init() {

	R.GET("/debug/pprof/", pprof.Index)
	R.GET("/debug/pprof/cmdline", pprof.Cmdline)
	R.GET("/debug/pprof/profile", pprof.Profile)
	R.GET("/debug/pprof/symbol", pprof.Symbol)
	R.GET("/debug/pprof/trace", pprof.Trace)

	R.GET("/h/merger", model.HandlerMerger)
	R.GET("/h/schedule", model.HandleGetSchedule)
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
