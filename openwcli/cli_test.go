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
appkey = "8c511cb683041f3589419440fab0a7b7710907022b0d035baea9001d529ca72f"

# App Secret
appid = "b4b1962d415d4d30ec71b28769fda585"

# Log file path
logdir = "./testdata/logs/"

# Data directory, store keys, databases, backups
datadir = "./testdata/data/"

# Wallet Summary Period
summaryperiod = "10s"

`

	c, err := config.NewConfigData("ini", []byte(confFile))
	if err != nil {
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