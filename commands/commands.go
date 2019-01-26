package commands

import (
	"github.com/blocktree/OpenWallet/cmd/utils"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/go-openw-cli/openwcli"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 通信节点命令
	Commands = []cli.Command{
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
			Usage:    "Show all wallet information",
			Action:   listwallet,
			Category: "WALLET COMMANDS",
			Flags: []cli.Flag{
				utils.SymbolFlag,
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
				utils.SymbolFlag,
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
