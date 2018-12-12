package openwcli

import (
	"github.com/blocktree/OpenWallet/log"
	"path/filepath"
	"testing"
)

func getTestOpenwCLI() *CLI {
	c := &Config{
		remoteserver: "127.0.0.1:8090",
		version: "1.0.0",
		appkey: "1234qwer",
		appsecret: "1234qwer",
		datadir: filepath.Join("test_openwcli", "data"),
		logdir: filepath.Join("test_openwcli", "log"),
		summaryperiod: "10s",
	}

	cli, err := NewCLI(c)
	if err != nil {
		log.Error("getTestOpenwCLI error:", err)
		return nil
	}
	return cli
}

func TestCLI_GenKeychain(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	err := cli.GenKeychain()
	if err != nil {
		log.Error("GenKeychain error:", err)
		return
	}
}
