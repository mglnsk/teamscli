package main

import (
	"log"

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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("ERROR reading config file ./config.yaml:")
		log.Fatal(err)
	}

	ctx := kong.Parse(&CLI)
	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
