package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "broker-gateway",
	Short: "Matthiola",
	Long:  "The broker gateway of project Matthiola, a distributed commodities OTC electronic trading system",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/config.toml", "config file")
	rootCmd.AddCommand(receiverCmd)
	rootCmd.AddCommand(senderCmd)

	viper.SetDefault("author", "rudeigerc <rudeigerc@gmail.com>")
	viper.SetDefault("license", "MIT")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

}
