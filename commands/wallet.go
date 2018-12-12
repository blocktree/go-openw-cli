package commands

import (
	"github.com/blocktree/OpenWallet/cmd/utils"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/go-openw-cli/openwcli"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 钱包命令
	CmdWallet = cli.Command{
		Name:      "wallet",
		Usage:     "Manage multi currency wallet",
		ArgsUsage: "",
		Category:  "Application COMMANDS",
		Description: `
You create, import, restore wallet

`,

		Subcommands: []cli.Command{
			{
				//获取钱包列表信息
				Name:     "list",
				Usage:    "Get all wallet information",
				Action:   newwallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.SymbolFlag,
				},
			},
			{
				//创建钱包
				Name:      "newwallet",
				Usage:     "new a currency wallet",
				ArgsUsage: "<symbol>",
				Action:    newwallet,
				Category:  "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.SymbolFlag,
				},
				Description: `
	wmd wallet new -s <symbol>

This command will start the wallet node, and create new wallet.

	`,
			},
		},
	}
)

//walletConfig 钱包配置
func newwallet(c *cli.Context) error {

	var (
		err error
	)

	conf := c.GlobalString("conf")
	config, err := openwcli.LoadConfig(conf)
	if err != nil {
		log.Error("unexpected error: ", err)
		return err
	}

	cli, err := openwcli.NewCLI(config)
	if err != nil {
		log.Error("unexpected error: ", err)
		return err
	}

	err = cli.NewWalletFlow()
	if err != nil {
		log.Error("unexpected error: ", err)
		return err
	}

	return nil
}
