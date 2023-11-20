package cli

import (
	"os"
	"strings"

	"github.com/k0kubun/go-ansi"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/collector"
	"github.com/voidint/g/version"
)

const (
	stableChannel   = "stable"
	unstableChannel = "unstable"
	archivedChannel = "archived"
)

func listRemote(ctx *cli.Context) (err error) {
	channel := ctx.Args().First()
	if channel != "" && channel != stableChannel && channel != unstableChannel && channel != archivedChannel {
		return cli.ShowSubcommandHelp(ctx)
	}

	c, err := collector.NewCollector(strings.Split(os.Getenv(mirrorEnv), ",")...)
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	var vs []*version.Version
	switch channel {
	case stableChannel:
		vs, err = c.StableVersions()
	case unstableChannel:
		vs, err = c.UnstableVersions()
	case archivedChannel:
		vs, err = c.ArchivedVersions()
	default:
		vs, err = c.AllVersions()
	}
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	var renderMode uint8
	switch ctx.String("output") {
	case "json":
		renderMode = jsonMode
	default:
		renderMode = rawMode
	}

	render(renderMode, installed(), vs, ansi.NewAnsiStdout())
	return nil
}
