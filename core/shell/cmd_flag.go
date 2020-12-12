///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package shell

import (
	"digger/agents"
	"digger/common"
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/urfave/cli"
	"os"
)

// Parse parses command flags using `github.com/urfave/cli`
func Parse(arguments []string) {
	appFlag := cli.NewApp()
	appFlag.Version = common.VERSION
	appFlag.HideVersion = true
	appFlag.Name = "digger"
	appFlag.Usage = "digger"
	appFlag.HelpName = "digger"
	// config file location
	appFlag.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "version, v",
			Usage:       `show version`,
			Destination: &showVersion,
		},
	}

	appFlag.Commands = []cli.Command{
		{
			Name:  "manager",
			Usage: "start as manager",
			Action: func(c *cli.Context) error {
				bootMode = common.ROLE_MANAGER
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "secret, s",
					Value:       "123456",
					Usage:       `custom global secret`,
					Destination: &secret,
				},
				cli.IntFlag{
					Name:        "port, p",
					Value:       9012,
					Usage:       "server http port",
					Destination: &port,
				},
				cli.StringFlag{
					Name:  "log-level, l",
					Value: "",
					Usage: `set log level, available options:
	(trace|debug|info|warn|error|fatal)`,
					Destination: &logLevel,
				},
				cli.StringFlag{
					Name:        "log-dir",
					Value:       "/var/log/digger",
					Usage:       "log directory",
					Destination: &logDir,
				},
				cli.StringFlag{
					Name:  "database, d",
					Value: "",
					Usage: `postgres connection string, format:
	postgres://<user>:<password>@<host>:<port>/<db>?sslmode=disable`,
					Destination: &dbConn,
				},
				cli.StringFlag{
					Name:  "redis, r",
					Value: "@127.0.0.1#1",
					Usage: `redis connection string, format:
	<password>@<host>:<port>#<db>`,
					Destination: &redisConn,
				},
				cli.StringFlag{
					Name:  "labels",
					Value: "",
					Usage: `node labels, format:
	key1=value1<,key2=value2>`,
					Destination: &labels,
				},
				cli.StringFlag{
					Name:        "ui-dir",
					Value:       "/var/www/html",
					Usage:       "ui dir",
					Destination: &uiDir,
				},
			},
		},
		{
			Name:  "worker",
			Usage: "start as worker",
			Action: func(c *cli.Context) error {
				bootMode = common.ROLE_WORKER
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "id",
					Value:       1,
					Usage:       `unique instance id`,
					Destination: &instanceId,
				},
				cli.StringFlag{
					Name:  "log-level",
					Value: "",
					Usage: `set log level, available options:
	(trace|debug|info|warn|error|fatal)`,
					Destination: &logLevel,
				},
				cli.StringFlag{
					Name:        "log-dir",
					Value:       "/var/log/digger",
					Usage:       "log directory",
					Destination: &logDir,
				},
				cli.StringFlag{
					Name:  "manager, m",
					Value: "localhost:9012",
					Usage: `manager address, format:
	<host>:<port>`,
					Destination: &managerUrl,
				},
				cli.StringFlag{
					Name:  "labels",
					Value: "",
					Usage: `node labels, format:
	key1=value1<,key2=value2>`,
					Destination: &labels,
				},
			},
		},
	}

	cli.AppHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .VisibleCommands}}

Commands:{{range .VisibleCategories}}
{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

Options:

   {{range $index, $option := .VisibleFlags}}{{if $index}}{{end}}{{$option}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

{{.Usage}}{{if .VisibleFlags}}

Options:

   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	cli.SubcommandHelpTemplate = `
Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

{{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

Commands:
{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{if .VisibleFlags}}

Options:

   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	appFlag.Action = func(c *cli.Context) error {
		if showVersion {
			cli.ShowVersion(c)
			os.Exit(0)
			return nil
		}
		cli.ShowAppHelp(c)
		os.Exit(0)
		return nil
	}

	err := appFlag.Run(arguments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 从环境覆盖配置
	utils.ExchangeEnvValue("DIGGER_SECRET", func(envValue string) {
		secret = envValue
	})
	utils.ExchangeEnvValue("DIGGER_PORT", func(envValue string) {
		p, err := convert.StrToInt(envValue)
		if err != nil {
			logger.Fatal("非法端口号: \"", envValue, "\": ", err)
		}
		port = p
	})
	utils.ExchangeEnvValue("DIGGER_LOG_LEVEL", func(envValue string) {
		logLevel = envValue
	})
	utils.ExchangeEnvValue("DIGGER_LOG_DIR", func(envValue string) {
		logDir = envValue
	})
	utils.ExchangeEnvValue("DIGGER_DATABASE", func(envValue string) {
		dbConn = envValue
	})
	utils.ExchangeEnvValue("DIGGER_MANAGER", func(envValue string) {
		managerUrl = envValue
	})
	utils.ExchangeEnvValue("DIGGER_REDIS", func(envValue string) {
		redisConn = envValue
	})
	utils.ExchangeEnvValue("DIGGER_LABELS", func(envValue string) {
		labels = envValue
	})
	utils.ExchangeEnvValue("DIGGER_UI_DIR", func(envValue string) {
		uiDir = envValue
	})
	utils.ExchangeEnvValue("DIGGER_ID", func(envValue string) {
		id, err := convert.StrToInt(envValue)
		if err != nil {
			logger.Fatal("无效ID: \"", envValue, "\": ", err)
		}
		instanceId = id
	})

	common.LogDir = logDir

	if bootMode == common.ROLE_NONE {
		os.Exit(0)
	} else if bootMode == common.ROLE_MANAGER {
		agents.Start(models.BootstrapConfig{
			BootMode:    common.ROLE_MANAGER,
			Port:        port,
			LogDir:      logDir,
			LogLevel:    logLevel,
			Secret:      secret,
			ManagerUrl:  managerUrl,
			DBString:    dbConn,
			RedisString: redisConn,
			InstanceId:  instanceId,
			Labels:      utils.ParseLabels(labels),
			UIDir:       uiDir,
		})
	} else if bootMode == common.ROLE_WORKER {
		agents.Start(models.BootstrapConfig{
			BootMode:    common.ROLE_WORKER,
			Port:        port,
			LogDir:      logDir,
			LogLevel:    logLevel,
			Secret:      secret,
			ManagerUrl:  managerUrl,
			DBString:    dbConn,
			RedisString: redisConn,
			InstanceId:  instanceId,
			Labels:      utils.ParseLabels(labels),
		})
	}
}
