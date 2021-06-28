package commands

import (
	"github.com/blocktree/go-openw-cli/v2/openwcli"
	"github.com/blocktree/openwallet/v2/log"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 通信节点命令
	Commands = []cli.Command{
		CmdVersion,
		{
			//生成keychain
			Name:      "genkeychain",
			Usage:     "Generate new keychain and print it",
			ArgsUsage: "",
			Action:    genkeychain,
			Category:  "OPENW-CLI COMMANDS",
		},
		{
			//登记节点
			Name:      "noderegister",
			Usage:     "create new keychain and register node to openw-server",
			ArgsUsage: "",
			Action:    noderegister,
			Category:  "OPENW-CLI COMMANDS",
		},
		{
			//节点信息
			Name:      "nodeinfo",
			Usage:     "show node information",
			ArgsUsage: "",
			Action:    nodeinfo,
			Category:  "OPENW-CLI COMMANDS",
		},
		{
			//获取钱包列表信息
			Name:     "listwallet",
			Usage:    "show all wallet information",
			Action:   listwallet,
			Category: "WALLET COMMANDS",
			Flags:    []cli.Flag{},
		},
		{
			//创建钱包
			Name:      "newwallet",
			Usage:     "create a new wallet",
			ArgsUsage: "<symbol>",
			Action:    newwallet,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "newaccount",
			Usage:     "create a new assets account",
			ArgsUsage: "<symbol>",
			Action:    newaccount,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listaccount",
			Usage:     "show all assets account",
			ArgsUsage: "<symbol>",
			Action:    listaccount,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "newaddress",
			Usage:     "select account to create batch address",
			ArgsUsage: "<symbol>",
			Action:    newaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "searchaddress",
			Usage:     "search address info",
			ArgsUsage: "<symbol>",
			Action:    searchaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "transfer",
			Usage:     "transfer certain amount of coins/tokens to destination address",
			ArgsUsage: "<symbol>",
			Action:    transfer,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "transferall",
			Usage:     "transfer all of coins/token to destination address",
			ArgsUsage: "<symbol>",
			Action:    transferall,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listsuminfo",
			Usage:     "show assets account summary info",
			ArgsUsage: "<symbol>",
			Action:    listsuminfo,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "setsum",
			Usage:     "setup assets account summary info",
			ArgsUsage: "<symbol>",
			Action:    setsum,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "startsum",
			Usage:     "start summary account task",
			ArgsUsage: "<symbol>",
			Action:    startsum,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
				FileFlag,
			},
		},
		{

			Name:      "updateinfo",
			Usage:     "update info from openw-server",
			ArgsUsage: "<symbol>",
			Action:    updateinfo,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listsymbol",
			Usage:     "show all symbols info",
			ArgsUsage: "<symbol>",
			Action:    listsymbol,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listtokencontract",
			Usage:     "show all token contract info",
			ArgsUsage: "<symbol>",
			Action:    listtokencontract,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listaddress",
			Usage:     "select account to show all address",
			ArgsUsage: "<symbol>",
			Action:    listaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "trustserver",
			Usage:     "start trusteeship wallet service for transmit node",
			ArgsUsage: "<symbol>",
			Action:    trustserver,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listtokenbalance",
			Usage:     "show account all token balance",
			ArgsUsage: "<symbol>",
			Action:    listtokenbalance,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "addtrustaddress",
			Usage:     "add trust address",
			ArgsUsage: "<symbol>",
			Action:    addtrustaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "listtrustaddress",
			Usage:     "show trust address list",
			ArgsUsage: "<symbol>",
			Action:    listtrustaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "enabletrustaddress",
			Usage:     "enable trust address",
			ArgsUsage: "<symbol>",
			Action:    enabletrustaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "disabletrustaddress",
			Usage:     "disable trust address",
			ArgsUsage: "<symbol>",
			Action:    disabletrustaddress,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "callabi",
			Usage:     "executes a new message call immediately without creating a transaction on the block chain.",
			ArgsUsage: "<symbol>",
			Action:    callabi,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "triggerabi",
			Usage:     "creates new transaction trigger smart contract on the block chain.",
			ArgsUsage: "<symbol>",
			Action:    triggerabi,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
		{

			Name:      "signhash",
			Usage:     "Select the private key of a specific address to sign the hash.",
			ArgsUsage: "<symbol>",
			Action:    signhash,
			Category:  "WALLET COMMANDS",
			Flags:     []cli.Flag{},
		},
	}
)

func getCLI(c *cli.Context) *openwcli.CLI {
	var (
		err error
	)

	conf := c.GlobalString("conf")
	config, err := openwcli.LoadConfig(conf)
	if err != nil {
		log.Error("unexpected error: ", err)
		return nil
	}

	cli, err := openwcli.NewCLI(config)
	if err != nil {
		log.Error("unexpected error: ", err)
		return nil
	}

	return cli
}

//register 注册
func noderegister(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.RegisterFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}
	return nil
}

//nodeinfo
func nodeinfo(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.GetNodeInfoFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//newwallet 创建钱包
func newwallet(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.NewWalletFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listwallet 钱包配置
func listwallet(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListWalletFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//newaccount 创建账户
func newaccount(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.NewAccountFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listaccount 账户列表
func listaccount(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListAccountFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//newaddress 创建地址
func newaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.NewAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//searchaddress 查询地址
func searchaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.SearchAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//transfer 转账交易
func transfer(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.TransferFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//transferall 转账全部资产交易
func transferall(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.TransferAllFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//setsum 设置汇总
func setsum(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.SetSumFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//startsum 定时汇总
func startsum(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {

		file := c.String("file")
		err := cli.StartSumFlow(file)
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//updateinfo 更新区块链资料库
func updateinfo(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.UpdateInfoFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listsymbol 查看所有币种信息列表
func listsymbol(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListSymbolFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listtokencontract 查看某个区块链所有代币合约信息列表
func listtokencontract(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListTokenContractFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listaddress 查看账户所有地址
func listaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//trustserver 启动后台托管钱包服务
func trustserver(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {

		err := cli.StartTrustServerFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

func genkeychain(c *cli.Context) error {

	err := openwcli.GenKeychainFlow()
	if err != nil {
		log.Error("unexpected error: ", err)
		return err
	}

	return nil
}

func listsuminfo(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListSumInfoFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//listtokenbalance
func listtokenbalance(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListTokenBalanceFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

// addtrustaddress
func addtrustaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.AddTrustAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

// listtrustaddress
func listtrustaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.ListTrustAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

// enabletrustaddress
func enabletrustaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.EnableTrustAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

// disabletrustaddress
func disabletrustaddress(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.DisableTrustAddressFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//callabi 直接调用ABI
func callabi(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.CallABIFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//triggerabi 触发ABI上链交易
func triggerabi(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.TriggerABIFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}

//signhash 消息签名
func signhash(c *cli.Context) error {

	if cli := getCLI(c); cli != nil {
		err := cli.SignHashFlow()
		if err != nil {
			log.Error("unexpected error: ", err)
			return err
		}
	}

	return nil
}
