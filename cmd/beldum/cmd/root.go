package cmd

import (
	"github.com/nanoteck137/beldum"
	"github.com/nanoteck137/beldum/config"
	"github.com/nanoteck137/beldum/core/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     beldum.AppName,
	Version: beldum.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(beldum.VersionTemplate(beldum.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
