package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/mglnsk/teamscli/cmd"
	"github.com/spf13/viper"
)

var CLI struct {
	Update   cmd.UpdateCmd   `cmd:"" help:"Update a channel."`
	Extract  cmd.ExtractCmd  `cmd:"" help:"Extract data from channel json"`
	Download cmd.DownloadCmd `cmd:"" help:"Download something"`
	Refresh  cmd.RefreshCmd  `cmd:"" help:"Just refresh the tokens"`
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	ctx := kong.Parse(&CLI)
	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
