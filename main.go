package main

import (
	"github.com/mglnsk/teamscli/cmd"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Update   cmd.UpdateCmd   `cmd:"" help:"Update a channel."`
	Extract  cmd.ExtractCmd  `cmd:"" help:"Extract data from channel json"`
	Download cmd.DownloadCmd `cmd:"" help:"Download something"`
	Refresh  cmd.RefreshCmd  `cmd:"" help:"Just refresh the tokens"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
