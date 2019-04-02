package main

import (
	"fmt"
	"github.com/blocktree/go-openw-cli/commands"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"sort"
)

const (
	VersionMajor = 1           // Major version component of the current release
	VersionMinor = 2           // Minor version component of the current release
	VersionPatch = 0           // Patch version component of the current release
	VersionMeta  = "stable"    // Version metadata to append to the version string
	Identifier   = "openw-cli" // Client identifier to advertise over the network
)

// Version holds the textual version string.
var Version = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

var (
	CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
	{{range .cmd.Subcommands}}{{.cmd.Name}}{{with .cmd.ShortName}}, {{.cmd}}{{end}}{{ "\t" }}{{.cmd.Usage}}
	{{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`

	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = NewApp(gitCommit, "the Wallet Manager Driver command line interface")
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = CommandHelpTemplate
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	app.Version = Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

func init() {
	// Initialize the CLI app and start openw-cli
	app.Name = Identifier
	app.Action = openwcli
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018 The OpenWallet Authors"
	app.Version = Version
	app.Commands = commands.Commands
	app.Flags = []cli.Flag{
		commands.AppNameFlag,
		commands.LogDirFlag,
		commands.DebugFlag,
		commands.ConfFlag,
	}

	sort.Sort(cli.CommandsByName(app.Commands))
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//openwcli
func openwcli(ctx *cli.Context) error {
	return nil
}
