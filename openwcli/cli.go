package openwcli

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/OpenWallet/common/file"
	"github.com/blocktree/OpenWallet/console"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/coreos/bbolt"
	"path/filepath"
	"time"
)


type CLI struct {
	//工具配置
	config *Config
	db     *openwallet.StormDB
}

// 初始化工具
func NewCLI(c *Config) (*CLI, error) {

	if len(c.appkey) == 0 {
		return nil, fmt.Errorf("appkey is empty. ")
	}

	if len(c.appsecret) == 0 {
		return nil, fmt.Errorf("appsecret is empty. ")
	}

	if len(c.remoteserver) == 0 {
		return nil, fmt.Errorf("remoteserver is empty. ")
	}

	//建立文件夹
	file.MkdirAll(c.datadir)
	file.MkdirAll(c.logdir)

	dbfile := filepath.Join(c.datadir, c.appkey+".db")

	//加载数据
	db, err := openwallet.OpenStormDB(
		dbfile,
		storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}),
	)


	if err != nil {
		return nil, err
	}

	cli := &CLI{
		config: c,
		db:     db,
	}
	return cli, nil
}

//checkConfig 检查配置加载完
func (cli *CLI) checkConfig() error {

	if cli.config == nil {
		return fmt.Errorf("config is not loaded. ")
	}

	if cli.db == nil {
		return fmt.Errorf("database is not loaded. ")
	}
	return nil
}

//RegisterFlow 注册节点流程
func (cli *CLI) RegisterFlow() error {

	var (
		current string
		confirm bool
	)

	if check := cli.checkConfig(); check != nil {
		return check
	}

	err := cli.db.Get(CLIBucket, CurrentKeychainKey, &current)
	if len(current) > 0  {
		//已经存在，提示是否需要覆盖
		confirm, _ = console.Stdin.PromptConfirm("The keychain already exist, do you want to regenerate current keychain?")
	} else {
		confirm = true
	}

	if confirm {
		//生成keychain
		err = cli.GenKeychain()
		if err != nil {
			return err
		}
	}

	//登记节点
	err = cli.RegisterOpenwServer()
	if err != nil {
		return err
	}

	return nil
}

//GetNodeInfo 获取节点信息
func (cli *CLI) GetNodeInfoFlow() error {

	var current string
	err := cli.db.Get(CLIBucket, CurrentKeychainKey, &current)
	if err != nil {
		return fmt.Errorf("The keychain not exist, please register node first. ")
	}

	var keychain Keychain
	err = cli.db.One("NodeID", current, &keychain)
	if err != nil {
		return fmt.Errorf("The keychain not exist, please register node first. ")
	}

	printKeychain(&keychain)

	return nil
}

//printKeychain 打印证书钥匙串
func printKeychain(keychain *Keychain) {

	//打印证书信息
	log.Notice("--------------- PRIVATE KEY ---------------")
	log.Notice(keychain.PrivateKey)
	fmt.Println()
	log.Notice("--------------- PUBLIC KEY ---------------")
	log.Notice(keychain.PrivateKey)
	log.Notice("--------------- NODE ID ---------------")
	log.Notice(keychain.NodeID)
	fmt.Println()
}

//NewWalletFlow 创建钱包流程
func (cli *CLI) NewWalletFlow() error {
	//TODO: WIP
	name, _ := console.InputText("wallet name:", true)
	log.Info(name)
	return nil
}

//ListWalletFlow
func (cli *CLI) ListWalletFlow() error {
	//TODO: WIP
	return nil
}

//NewAccountFlow
func (cli *CLI) NewAccountFlow() error {
	//TODO: WIP
	return nil
}

//ListAccountFlow
func (cli *CLI) ListAccountFlow() error {
	//TODO: WIP
	return nil
}

//NewAddressFlow
func (cli *CLI) NewAddressFlow() error {
	//TODO: WIP
	return nil
}

//SearchAddressFlow
func (cli *CLI) SearchAddressFlow() error {
	//TODO: WIP
	return nil
}

//TransferFlow
func (cli *CLI) TransferFlow() error {
	//TODO: WIP
	return nil
}

//SetSumFlow
func (cli *CLI) SetSumFlow() error {
	//TODO: WIP
	return nil
}

//StartSumFlow
func (cli *CLI) StartSumFlow() error {
	//TODO: WIP
	return nil
}