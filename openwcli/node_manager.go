package openwcli

import (
	"fmt"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/owtp"
)

//GenKeychain 生成新密钥对
func (cli *CLI) GenKeychain() error {

	if check := cli.checkConfig(); check != nil {
		return check
	}

	//随机创建证书
	cert := owtp.NewRandomCertificate()
	if len(cert.PrivateKeyBytes()) == 0 {
		return fmt.Errorf("create keychain failed ")
	}

	keychain := NewKeychain(cert)

	//保存到数据库
	err := cli.db.Save(keychain)
	if err != nil {
		return fmt.Errorf("save new keychain failed. unexpected error: %v", err)
	}

	err = cli.db.Set(CLIBucket, CurrentKeychainKey, keychain.NodeID)
	if err != nil {
		return fmt.Errorf("update current keychain failed. unexpected error: %v", err)
	}

	log.Info("Create keychain successfully.")

	//打印密钥对
	printKeychain(keychain)

	return nil
}


//RegisterOpenwServer 注册节点到openw-server
func (cli *CLI) RegisterOpenwServer() error {

	//	TODO: 登记节点到openw-server
	return nil
}
