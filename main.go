package main

import (
	_ "embed"
	"log"
	"log/slog"
	"os"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/fatal"
	"github.com/richardwilkes/toolbox/log/tracelog"
	"github.com/richardwilkes/unison"
)

func main() {
	cmdline.AppName = "Gemiclip"
	cmdline.AppCmdName = "gemiclip"
	cmdline.CopyrightStartYear = "2024"
	cmdline.CopyrightHolder = "Harry Nguyen"
	cmdline.AppIdentifier = "com.gemiclip.app"

	unison.AttachConsole()

	cl := cmdline.New(true)
	cl.Parse(os.Args[1:])
	slog.SetDefault(slog.New(tracelog.New(log.Default().Writer(), slog.LevelInfo)))

	unison.Start(unison.StartupFinishedCallback(func() {
		_, err := NewMarkdownWindow(unison.PrimaryDisplay().Usable.Point)
		fatal.IfErr(err)
	})) // Never returns
}
