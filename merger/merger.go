package merger

import (
	"github.com/mx5566/mergernew/config"
	"github.com/mx5566/mergernew/route"
	"strconv"
)

func Run() error {

	var err error
	if config.Config.Port == 0 {
		err = route.GinR.Run(":5050")
	} else {
		err = route.GinR.Run(":" + strconv.Itoa(int(config.Config.Port)))
	}

	return err
}
