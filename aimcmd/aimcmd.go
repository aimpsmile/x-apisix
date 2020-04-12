package aimcmd

import (
	"github.com/micro-in-cn/x-apisix/core/aimerror"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"log"
)

const (
	TMPL = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`
)

var aime = aimerror.Errors([]error{})

//	获取错误列表是否
func AimErrors() aimerror.Errors {
	return aime
}

//	只有在逻辑初始化场景才可以使用,,追加错误
func AimErrorsAppend(newErrors ...error) {
	aime = AimErrors().Add(newErrors...)
}

type CommandI interface {
	Name() string
	Version() string
	Usage() string
	SubFlags() []string
	Commands(options ...micro.Option) []*cli.Command
}

// Setup sets up a cli.App
func setup(cmd CommandI, app *cli.App) {
	app.Name = cmd.Name()
	app.Version = cmd.Version()
	app.Usage = cmd.Usage()
	app.Flags = []cli.Flag{}
	app.EnableBashCompletion = true
	app.CustomAppHelpTemplate = TMPL
	app.Commands = append(app.Commands, cmd.Commands()...)
	app.OnUsageError = func(context *cli.Context, err error, isSubcommand bool) error {
		log.Println(context.Command.Name, "---->", err.Error())
		return nil
	}
	app.Action = cli.ShowAppHelp
}

func Init(cmd CommandI) {
	//	构建cli命令步骤
	app := cli.NewApp()
	setup(cmd, app)
	//	处理flag
	nflags := cmd.SubFlags()
	//	运行命令管理器
	if err := app.Run(nflags); err != nil {
		log.Panic(err)
	}
	//	批量错误处理
	if AimErrors().IsError() {
		log.Panic(AimErrors().Error())
	}
}
