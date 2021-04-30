package route

import (
	"net/http"
)

//实现IOdServer的接口，以及http提供ServeHttp方法
type OdServer struct {
	router MethodMaps
}

type IOdServer interface {
	GET(url string, f HandlerFunc)
	POST(url string, f HandlerFunc)
	PUT(url string, f HandlerFunc)
	DELETE(url string, f HandlerFunc)
}

type HandlerMapped struct {
	f HandlerFunc
}

//接口函数单位，即我们编写代码逻辑的函数
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

func Default() *OdServer {
	return &OdServer{
		router: NewRouter(),
	}
}

//实现Handler接口，匹配方法以及路径
func (o *OdServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//转发给doHandler进行执行
	o.doHandler(w, req)
}

//判断需要执行的Http Method，从而查找对应的接口并且执行
func (o *OdServer) doHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			if hm, ok := o.router.GetMapping(req.URL.Path); ok {
				//if hm, ok := o.router.GetMapping(req.URL.RequestURI()); ok {
				hm.f(w, req)
			}
		}
	case http.MethodPost:
		{
			if hm, ok := o.router.PostMapping(req.URL.Path); ok {
				hm.f(w, req)
			}
		}
	case http.MethodDelete:
		{
			if hm, ok := o.router.DeleteMapping(req.URL.String()); ok {
				hm.f(w, req)
			}
		}
	case http.MethodPut:
		{
			if hm, ok := o.router.PutMapping(req.URL.String()); ok {
				hm.f(w, req)
			}
		}
	default:
		{

		}
	}
}

func (o *OdServer) GET(url string, f HandlerFunc) {
	o.router.GetAdd(url, HandlerMapped{f: f})
}
func (o *OdServer) POST(url string, f HandlerFunc) {
	o.router.PostAdd(url, HandlerMapped{f: f})
}
func (o *OdServer) PUT(url string, f HandlerFunc) {
	o.router.PutAdd(url, HandlerMapped{f: f})
}
func (o *OdServer) DELETE(url string, f HandlerFunc) {
	o.router.DeleteAdd(url, HandlerMapped{f: f})
}
