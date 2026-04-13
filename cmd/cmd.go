package cmd

import (
	config2 "VPMBuilder/internal/config"
	content2 "VPMBuilder/internal/content"
	"VPMBuilder/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "vpm-builder",
	Short: "VPM Builder",
	Long:  "A Tool to build VPM packages",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		startup()
		return nil
	},
}

var configPath string
var config *config2.Config
var content *content2.Content

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.json", "config file path")
}

func startup() {
	log.Info("Initializing VPM Builder...")
	config = config2.NewConfig()
	err := config.Parse(configPath)
	if err != nil {
		log.ErrorE("Failed to parse config", zap.Error(err))
	}
	log.SetLogLevel(config.LogLevel)
	log.Info("Config loaded.")
	log.Info("Initializing Content...")
	content = content2.NewContent(config)
	log.Info("VPM Builder initialized.")
}

func Execute() error {
	printLogo()
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}
