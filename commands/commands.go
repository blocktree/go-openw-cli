package commands

import (
	"github.com/blocktree/openwallet/log"
	"github.com/blocktree/go-openw-cli/openwcli"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 通信节点命令
	Commands = []cli.Command{
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
			Flags: []cli.Flag{
			},
		},
		{
			//创建钱包
			Name:      "newwallet",
			Usage:     "create a new wallet",
			ArgsUsage: "<symbol>",
			Action:    newwallet,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "newaccount",
			Usage:     "create a new assets account",
			ArgsUsage: "<symbol>",
			Action:    newaccount,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "listaccount",
			Usage:     "show all assets account",
			ArgsUsage: "<symbol>",
			Action:    listaccount,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "newaddress",
			Usage:     "select account to create batch address",
			ArgsUsage: "<symbol>",
			Action:    newaddress,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "searchaddress",
			Usage:     "search address info",
			ArgsUsage: "<symbol>",
			Action:    searchaddress,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "transfer",
			Usage:     "create a transaction",
			ArgsUsage: "<symbol>",
			Action:    transfer,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "setsum",
			Usage:     "setup assets account summary info",
			ArgsUsage: "<symbol>",
			Action:    setsum,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
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
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "listsymbol",
			Usage:     "show all symbols info",
			ArgsUsage: "<symbol>",
			Action:    listsymbol,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "listtokencontract",
			Usage:     "show all token contract info",
			ArgsUsage: "<symbol>",
			Action:    listtokencontract,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "listaddress",
			Usage:     "select account to show all address",
			ArgsUsage: "<symbol>",
			Action:    listaddress,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
		},
		{

			Name:      "trustserver",
			Usage:     "start trusteeship wallet service for transmit node",
			ArgsUsage: "<symbol>",
			Action:    trustserver,
			Category:  "WALLET COMMANDS",
			Flags: []cli.Flag{
			},
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