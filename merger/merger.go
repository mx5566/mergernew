package merger

import (
	"github.com/mx5566/mergernew/route"
	"net/http"
)

func Run() {
	_ = http.ListenAndServe(":5050", route.R)
}
