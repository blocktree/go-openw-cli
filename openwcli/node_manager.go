package openwcli

import (
	"fmt"
	"github.com/blocktree/openwallet/v2/owtp"
)

//GenKeychain 生成新密钥对
func GenKeychain() (*Keychain, error) {

	//if check := cli.checkConfig(); check != nil {
	//	return nil, check
	//}

	//随机创建证书
	cert := owtp.NewRandomCertificate()
	if len(cert.PrivateKeyBytes()) == 0 {
		return nil, fmt.Errorf("create keychain failed ")
	}

	keychain := NewKeychain(cert)

	return keychain, nil

	////保存到数据库
	//err := cli.db.Save(keychain)
	//if err != nil {
	//	return nil, fmt.Errorf("save new keychain failed. unexpected error: %v", err)
	//}
	//
	//err = cli.db.Set(CLIBucket, CurrentKeychainKey, keychain.NodeID)
	//if err != nil {
	//	return nil, fmt.Errorf("update current keychain failed. unexpected error: %v", err)
	//}

	return keychain, nil
}

//SaveCurrentKeychain 保存新密钥对到本地缓存
func (cli *CLI) SaveCurrentKeychain(keychain *Keychain) error {

	if check := cli.checkConfig(); check != nil {
		return check
	}

	_, err := cli.getDB()
	if err != nil {
		return err
	}
	defer cli.closeDB()

	//保存到数据库
	err = cli.db.Save(keychain)
	if err != nil {
		return fmt.Errorf("save new keychain failed. unexpected error: %v", err)
	}

	err = cli.db.Set(CLIBucket, CurrentKeychainKey, keychain.NodeID)
	if err != nil {
		return fmt.Errorf("update current keychain failed. unexpected error: %v", err)
	}

	return nil
}

//GetKeychain
func (cli *CLI) GetKeychain() (*Keychain, error) {

	_, err := cli.getDB()
	if err != nil {
		return nil, err
	}
	defer cli.closeDB()

	var current string
	err = cli.db.Get(CLIBucket, CurrentKeychainKey, &current)
	if err != nil {
		return nil, fmt.Errorf("The keychain not exist, please register node first. ")
	}

	var keychain Keychain
	err = cli.db.One("NodeID", current, &keychain)
	if err != nil {
		return nil, fmt.Errorf("The keychain not exist, please register node first. ")
	}

	return &keychain, nil
}

//RegisterOnServer 注册节点到openw-server
func (cli *CLI) RegisterOnServer() error {

	//登记节点到openw-server
	return cli.api.BindAppDevice()
}
