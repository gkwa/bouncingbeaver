package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose int
)

var rootCmd = &cobra.Command{
	Use:   "bouncingbeaver",
	Short: "A tool for processing DynamoDB data",
	Long:  "A command-line tool that demonstrates unmarshaling DynamoDB AttributeValue format",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bouncingbeaver.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "verbose output (can be used multiple times)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bouncingbeaver")
	}
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
