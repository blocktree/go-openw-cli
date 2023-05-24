package openwcli

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/go-openw-sdk/v2/openwsdk"
	"github.com/blocktree/openwallet/v2/hdkeystore"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/owtp"
	"path/filepath"
	"testing"
)

func init() {
	owtp.Debug = false
}

func getTestOpenwCLI() *CLI {

	confFile := filepath.Join("conf", "test.ini")
	//confFile := filepath.Join("conf", "prod.ini")

	c, err := config.NewConfig("ini", confFile)
	if err != nil {
		log.Error("NewConfigData error:", err)
		return nil
	}
	conf := NewConfig(c)
	cli, err := NewCLI(conf)
	if err != nil {
		log.Error("getTestOpenwCLI error:", err)
		//return nil
	}
	return cli

}

func TestChangePwd(t *testing.T) {
	cli := getTestOpenwCLI()
	if cli == nil {
		return
	}

	keystore := hdkeystore.NewHDKeystore(
		cli.config.keydir,
		hdkeystore.StandardScryptN,
		hdkeystore.StandardScryptP,
	)

	//随机生成keystore
	key, filePath, err := hdkeystore.StoreHDKey(
		cli.config.keydir,
		"pwd",
		"12345678",
		hdkeystore.StandardScryptN,
		hdkeystore.StandardScryptP,
	)
	if err != nil {
		log.Error("StoreHDKey error:", err)
		return
	}

	wallet := &openwsdk.Wallet{
		Alias:    "pwd",
		WalletID: key.KeyID,
	}

	//用新密码加密
	err = keystore.StoreKey(filePath, key, "1234qwer")
	if err != nil {
		log.Error("StoreKey error:", err)
		return
	}

	// 解密钱包
	_, err = cli.getLocalKeyByWallet(wallet, "1234qwer")
	if err != nil {
		log.Error("can not use new password unlock wallet:", err)
		return
	}
}
