package commands

import (
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/go-openw-cli/openwcli"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 通信节点命令
	CmdNode = cli.Command{
		Name:      "node",
		Usage:     "Manage openw node",
		ArgsUsage: "",
		Category:  "OPENW-CLI COMMANDS",
		Description: `
Use node commands to register openw-server

`,
		Subcommands: []cli.Command{
			{
				//登记节点
				Name:      "register",
				Usage:     "create new keychain and register node to openw-server",
				ArgsUsage: "<init>",
				Action:    register,
				Category:  "NODE COMMANDS",
			},
		},
	}
)

//register 注册
func register(c *cli.Context) error {

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

	err = cli.RegisterFlow()
	if err != nil {
		log.Error("unexpected error: ", err)
		return err
	}

	return nil
}
