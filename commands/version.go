package commands

import (
	"fmt"
	"github.com/blocktree/go-openw-cli/openwcli"
	"gopkg.in/urfave/cli.v1"
)

var (
	// 钱包命令
	CmdVersion = cli.Command{
		Name:      "version",
		Usage:     "show version information",
		ArgsUsage: "",
		Action:    version,
		Category:  "OPENW-CLI COMMANDS",
	}
)

//walletConfig 钱包配置
func version(c *cli.Context) error {
	fmt.Printf("FixAppID: %s\n", openwcli.FixAppID)
	fmt.Printf("Version: %s\n", openwcli.Version)
	fmt.Printf("GitRev: %s\n", openwcli.GitRev)
	fmt.Printf("BuildTime: %s\n", openwcli.BuildTime)
	return nil
}
