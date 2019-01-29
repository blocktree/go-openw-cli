package openwcli

import (
	"encoding/json"
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/OpenWallet/console"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/OpenWallet/timer"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/coreos/bbolt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxAddresNum = 2000
)

type CLI struct {
	config      *Config             //工具配置
	db          *openwallet.StormDB //本地数据库
	api         *openwsdk.APINode   //api
	summaryTask *SummaryTask        //汇总任务
}

func init() {
	owtp.Debug = false
}

// 初始化工具
func NewCLI(c *Config) (*CLI, error) {

	if len(c.appkey) == 0 {
		return nil, fmt.Errorf("appkey is empty. ")
	}

	if len(c.appid) == 0 {
		return nil, fmt.Errorf("appid is empty. ")
	}

	if len(c.remoteserver) == 0 {
		return nil, fmt.Errorf("remoteserver is empty. ")
	}

	dbfile := filepath.Join(c.dbdir, c.appid+".db")

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

	//配置日志
	SetupLog(c.logdir, "openwcli.log", false)

	keychain, err := cli.GetKeychain()
	if keychain != nil {
		cli.setupAPISDK(keychain)
	}

	return cli, nil
}

//setupAPI 配置APISDK
func (cli *CLI) setupAPISDK(keychain *Keychain) error {

	if keychain != nil {
		cert, _ := keychain.Certificate()
		sdkConfig := &openwsdk.APINodeConfig{
			AppID:              cli.config.appid,
			AppKey:             cli.config.appkey,
			ConnectType:        owtp.HTTP,
			Host:               cli.config.remoteserver,
			EnableSignature:    false,
			EnableKeyAgreement: false,
			Cert:               cert,
		}

		apiSDK := openwsdk.NewAPINode(sdkConfig)
		cli.api = apiSDK
	}

	return nil
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
		confirm  bool
		keychain *Keychain
	)

	if check := cli.checkConfig(); check != nil {
		return check
	}

	keychain, err := cli.GetKeychain()
	if keychain != nil {
		//已经存在，提示是否需要覆盖
		confirm, _ = console.Stdin.PromptConfirm("The keychain already exist, do you want to regenerate current keychain?")
	} else {
		confirm = true
	}

	if confirm {
		//生成keychain
		keychain, err = cli.GenKeychain()
		if err != nil {
			return err
		}

		log.Info("Create keychain successfully.")

		//打印密钥对
		printKeychain(keychain)
	}

	//配置APISDK
	err = cli.setupAPISDK(keychain)
	if err != nil {
		return err
	}

	//登记节点
	err = cli.RegisterOnServer()
	if err != nil {
		return err
	}

	log.Info("Register node on opew-server successfully.")

	return nil
}

//GetNodeInfo 获取节点信息
func (cli *CLI) GetNodeInfoFlow() error {

	keychain, err := cli.GetKeychain()
	if err != nil {
		return err
	}

	printKeychain(keychain)

	return nil
}

//printKeychain 打印证书钥匙串
func printKeychain(keychain *Keychain) {
	//打印证书信息
	log.Notice("--------------- PRIVATE KEY ---------------")
	log.Notice(keychain.PrivateKey)
	log.Notice("--------------- PUBLIC KEY ---------------")
	log.Notice(keychain.PublicKey)
	log.Notice("--------------- NODE ID ---------------")
	log.Notice(keychain.NodeID)
}

//NewWalletFlow 创建钱包流程
func (cli *CLI) NewWalletFlow() error {

	var (
		password string
		name     string
		err      error
	)

	if cli.api == nil {
		return err
	}

	// 等待用户输入钱包名字
	name, err = console.InputText("Enter wallet's name: ", true)

	// 等待用户输入密码
	password, err = console.InputPassword(false, 3)

	err = cli.CreateWalletOnServer(name, password)
	if err != nil {
		return err
	}

	return nil
}

//ListWalletFlow
func (cli *CLI) ListWalletFlow() error {
	wallets, _ := cli.GetWalletsOnServer()
	cli.printWalletList(wallets)
	return nil
}

//NewAccountFlow
func (cli *CLI) NewAccountFlow() error {

	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	//:输入钱包密码
	// 等待用户输入密码
	password, err := console.InputPassword(false, 3)
	if err != nil {
		return err
	}

	//:输入账户别名
	// 等待用户输入钱包名字
	name, err := console.InputText("Enter account's name: ", true)
	if err != nil {
		return err
	}

	//:输入币种类别
	// 等待用户输入钱包名字
	symbol, err := console.InputText("Enter account's symbol: ", true)
	if err != nil {
		return err
	}

	//创建新账户
	err = cli.CreateAccountOnServer(name, password, symbol, wallet)
	if err != nil {
		return err
	}

	return nil
}

//ListAccountFlow
func (cli *CLI) ListAccountFlow() error {

	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	accounts, _ := cli.GetAccountsOnServer(wallet.WalletID)
	cli.printAccountList(accounts)
	return nil
}

//NewAddressFlow
func (cli *CLI) NewAddressFlow() error {

	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	//:选择账户
	account, err := cli.selectAccountStep(wallet.WalletID)
	if err != nil {
		return err
	}

	// 输入地址数量
	count, err := console.InputNumber("Enter the number of addresses you want: ", false)
	if err != nil {
		return err
	}

	if count > maxAddresNum {
		return fmt.Errorf("The number of addresses can not exceed %d ", maxAddresNum)
	}

	err = cli.CreateAddressOnServer(account.WalletID, account.AccountID, count)
	if err != nil {
		return err
	}
	return nil
}

//SearchAddressFlow
func (cli *CLI) SearchAddressFlow() error {

	var (
		password string
	)

	// 等待用户输入地址
	addr, err := console.InputText("Enter address: ", true)
	if err != nil {
		return err
	}

	//是否需要显示地址私钥，需要必须填入密码
	confirm, _ := console.Stdin.PromptConfirm("Do want to show address private key?")
	if confirm {
		// 等待用户输入密码
		password, err = console.InputPassword(false, 3)
		if err != nil {
			return err
		}
	}

	address, err := cli.SearchAddressOnServer(addr)
	if err != nil {
		return err
	}

	cli.printAddressList(address.WalletID, []*openwsdk.Address{address}, password)

	return nil
}

//TransferFlow
func (cli *CLI) TransferFlow() error {
	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	//:选择账户
	account, err := cli.selectAccountStep(wallet.WalletID)
	if err != nil {
		return err
	}

	// 等待用户输入合约地址
	contractAddress, err := console.InputText("Enter contract address: ", false)
	if err != nil {
		return err
	}
	// 等待用户输入接收地址
	to, err := console.InputText("Enter received address: ", true)
	if err != nil {
		return err
	}

	// 等待用户输入发送数量
	amount, err := console.InputRealNumber("Enter amount to send: ", true)
	if err != nil {
		return err
	}

	// 等待用户输入密码
	password, err := console.InputPassword(false, 3)
	if err != nil {
		return err
	}

	err = cli.Transfer(wallet, account, contractAddress, to, amount, password)
	if err != nil {
		return err
	}

	return nil
}

//SetSumFlow
func (cli *CLI) SetSumFlow() error {

	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	//:选择账户
	account, err := cli.selectAccountStep(wallet.WalletID)
	if err != nil {
		return err
	}

	sumAddress, err := console.InputText("Enter account's summary address: ", true)
	if err != nil {
		return err
	}

	threshold, err := console.InputText("Enter account's summary threshold: ", true)
	if err != nil {
		return err
	}

	minTransfer, err := console.InputText("Enter address's minimum transfer amount: ", true)
	if err != nil {
		return err
	}

	retainedBalance, err := console.InputText("Enter address's retained balance: ", true)
	if err != nil {
		return err
	}

	confirms, err := console.InputNumber("Enter how many confirms can transfer: ", true)
	if err != nil {
		return err
	}

	err = cli.SetSummaryInfo(account.WalletID, account.AccountID, sumAddress, threshold, minTransfer, retainedBalance, confirms)
	if err != nil {
		return err
	}

	log.Infof("setup summary info successfully")
	return nil
}

//StartSumFlow
func (cli *CLI) StartSumFlow(file string) error {

	var (
		endRunning  = make(chan bool, 1)
		manual      = true //手动选择
		summaryTask SummaryTask
		taskFile    string
	)

	cycleTime := cli.config.summaryperiod
	if len(cycleTime) == 0 {
		cycleTime = "1m"
	}

	cycleSec, err := time.ParseDuration(cycleTime)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		taskFile, err = console.InputText("Enter summary task json file path: ", false)
		if err != nil {
			return err
		}
	} else {
		taskFile = file
	}

	taskJSON, err := ioutil.ReadFile(taskFile)
	if err == nil {

		err = json.Unmarshal(taskJSON, &summaryTask)
		if err == nil {
			manual = false
		}
	}

	if manual {

		//:选择钱包
		wallet, selectErr := cli.selectWalletStep()
		if selectErr != nil {
			return selectErr
		}

		//:选择账户
		account, selectErr := cli.selectAccountStep(wallet.WalletID)
		if selectErr != nil {
			return selectErr
		}

		// 等待用户输入密码
		password, selectErr := console.InputPassword(false, 3)
		if selectErr != nil {
			return selectErr
		}

		summaryAccountTask := &SummaryAccountTask{
			AccountID: account.AccountID,
			Contracts: []string{},
		}

		summaryWalletTask := &SummaryWalletTask{
			WalletID: wallet.WalletID,
			Password: password,
			wallet:   wallet,
			Accounts: []*SummaryAccountTask{
				summaryAccountTask,
			},
		}

		summaryTask = SummaryTask{
			Wallets: []*SummaryWalletTask{
				summaryWalletTask,
			},
		}
	}

	cli.summaryTask = &summaryTask

	log.Infof("The timer for summary task start now. Execute by every %v seconds.", cycleSec.Seconds())

	//启动钱包汇总程序
	sumTimer := timer.NewTask(cycleSec, cli.SummaryTask)
	sumTimer.Start()

	<-endRunning

	return nil
}

//UpdateInfoFlow
func (cli *CLI) UpdateInfoFlow() error {

	err := cli.UpdateSymbols()
	if err != nil {
		return nil
	}

	log.Infof("update info successfully")

	return nil
}

//ListSymbolFlow
func (cli *CLI) ListSymbolFlow() error {
	list, err := cli.GetSymbolList()
	if err != nil {
		return err
	}
	cli.printSymbolList(list)
	return nil
}

//ListTokenContractFlow
func (cli *CLI) ListTokenContractFlow() error {

	symbol, err := console.InputText("Enter symbol: ", true)
	if err != nil {
		return err
	}
	symbol = strings.ToUpper(symbol)
	list, err := cli.GetTokenContractList("Symbol", symbol)
	if err != nil {
		return err
	}
	cli.printTokenContractList(list)
	return nil
}

//ListAddressFlow
func (cli *CLI) ListAddressFlow() error {

	//:选择钱包
	wallet, err := cli.selectWalletStep()
	if err != nil {
		return err
	}

	//:选择账户
	account, err := cli.selectAccountStep(wallet.WalletID)
	if err != nil {
		return err
	}

	offset, err := console.InputNumber("Enter offset: ", true)
	if err != nil {
		return err
	}

	limit, err := console.InputNumber("Enter limit: ", true)
	if err != nil {
		return err
	}

	addresses, err := cli.GetAddressesOnServer(account.WalletID, account.AccountID, int(offset), int(limit))
	if err != nil {
		return err
	}

	err = cli.printAddressList(account.WalletID, addresses, "")
	if err != nil {
		return err
	}

	return nil
}

//selectWalletStep 选择钱包操作
func (cli *CLI) selectWalletStep() (*openwsdk.Wallet, error) {

	wallets, _ := cli.GetWalletsOnServer()
	cli.printWalletList(wallets)

	fmt.Printf("[Please select a wallet] \n")

	//选择钱包
	num, err := console.InputNumber("Enter wallet No.: ", true)
	if err != nil {
		return nil, err
	}

	if int(num) >= len(wallets) {
		return nil, fmt.Errorf("Input number is out of index! ")
	}

	wallet := wallets[num]
	return wallet, nil
}

//selectAccountStep 选择资产账户操作
func (cli *CLI) selectAccountStep(walletID string) (*openwsdk.Account, error) {

	accounts, _ := cli.GetAccountsOnServer(walletID)
	cli.printAccountList(accounts)

	fmt.Printf("[Please select a account] \n")

	//选择钱包
	num, err := console.InputNumber("Enter account No.: ", true)
	if err != nil {
		return nil, err
	}

	if int(num) >= len(accounts) {
		return nil, fmt.Errorf("Input number is out of index! ")
	}

	account := accounts[num]
	return account, nil
}
