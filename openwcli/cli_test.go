package openwcli

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/OpenWallet/log"
)

func getTestOpenwCLI() *CLI {

	confFile := `
# Remote Server
remoteserver = "120.78.83.180"

# API Version
version = "1.0.0"

# App Key
appkey = "faa14b5e2cf119cd6d38bda45b49eb02b333a1b1ff6f10703acb554011ebfb1e"

# App Secret
appid = "8df7420d3917afa0172ea9c85e07ab55"

# Log file path
logdir = "./testdata/logs/"

# Data directory, store keys, databases, backups
datadir = "./testdata/data/"

# Wallet Summary Period
summaryperiod = "10s"

`

	c, err := config.NewConfigData("ini", []byte(confFile))
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