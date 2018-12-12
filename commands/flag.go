package commands

import "gopkg.in/urfave/cli.v1"

var (

	AppNameFlag = cli.StringFlag{
		Name: "name",
		Usage: "Application Name",
	}

	SymbolFlag = cli.StringFlag{
		Name: "symbol, s",
		Usage: "Currency symbol",
	}

	BatchFlag = cli.BoolFlag{
		Name: "batch",
		Usage: "Create address with batch",
	}

	InitFlag = cli.BoolFlag{
		Name: "i, init",
		Usage: "Init operate",
	}

	LogDirFlag = cli.StringFlag{
		Name: "logdir",
		Usage: "log files directory",
	}

	DebugFlag = cli.BoolFlag{
		Name: "debug",
		Usage: "print debug log info",
	}

	PathFlag = cli.StringFlag{
		Name: "path, p",
		Usage: "directory path",
	}

	FileFlag = cli.StringFlag{
		Name: "file, f",
		Usage: "file path",
	}

	ConfFlag = cli.StringFlag{
		Name: "conf, c",
		Usage: "config file path",
	}
)
