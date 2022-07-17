package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	port     int
	key      string
	postgres string
	rootCmd  = &cobra.Command{
		Use:   "twitter-mentions",
		Short: "monitor mentions of a certain account on twitter",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVar(&port, "port", 8081, "port to listen on")
	rootCmd.PersistentFlags().StringVar(&key, "key", "key.pem", "path to jwt private key file")
	rootCmd.PersistentFlags().String("postgres", "postgresql://twitter:abcde@127.0.0.1:5432/tweets", "postgres connection string")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("postgres", rootCmd.PersistentFlags().Lookup("postgres"))

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rulesCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
