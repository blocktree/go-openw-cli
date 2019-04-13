package openwcli

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/openwallet/owtp"
	"path/filepath"
)


func init() {
	owtp.Debug = false
}

func getTestOpenwCLI() *CLI {

	confFile := filepath.Join("conf", "node.ini")

	c, err := config.NewConfig("ini", confFile)
	if err != nil {
		log.Error("NewConfigData error:", err)
		return nil
	}
	conf := NewConfig(c)
	cli, err := NewCLI(conf)
	if err != nil {
		log.Error("getTestOpenwCLI error:", err)
		return nil
	}
	return cli
}