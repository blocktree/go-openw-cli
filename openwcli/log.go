package openwcli

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/blocktree/openwallet/common/file"
	"github.com/blocktree/openwallet/log"
	"path/filepath"
)

//SetupLog 配置日志
func SetupLog(logDir, logFile string, debug bool) {

	//记录日志
	logLevel := log.LevelInformational
	if debug {
		logLevel = log.LevelDebug
	}

	if len(logDir) > 0 {
		file.MkdirAll(logDir)
		logFile := filepath.Join(logDir, logFile)
		logConfig := fmt.Sprintf(`{"filename":"%s","level":%d,"daily":true,"maxdays":7,"maxsize":0}`, logFile, logLevel)
		//log.Println(logConfig)
		log.SetLogger(logs.AdapterFile, logConfig)
		log.SetLogger(logs.AdapterConsole, logConfig)
	} else {
		log.SetLevel(logLevel)
	}
}
