package openwcli

import (
	"github.com/blocktree/OpenWallet/log"
	"testing"
)

func TestCLI_GenKeychain(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	keychain, err := cli.GenKeychain()
	if err != nil {
		log.Error("GenKeychain error:", err)
		return
	}
	printKeychain(keychain)
}