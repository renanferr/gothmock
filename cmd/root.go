package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultPort    = ":6666"
	DefaultStatus  = http.StatusOK
	DefaultContent = "application/json"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	content     string
	status      int
	addr        string

	rootCmd = &cobra.Command{
		Use:   "gothmock [API specification] [filepath]",
		Short: "Mocks an API",
		Long:  `Mocks an API from an API specification file`,
		Args:  cobra.MinimumNArgs(1),
		// Run:   Mock,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "Renan Ferreira <https://github.com/renanferr>", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.PersistentFlags().StringVar(&addr, "port", DefaultPort, "port to run application server on")
	rootCmd.PersistentFlags().StringVar(&content, "content", DefaultContent, "response content-type")
	rootCmd.PersistentFlags().IntVar(&status, "status", DefaultStatus, "response status code")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Renan Ferreira <https://github.com/renanferr>")
	viper.SetDefault("license", "MIT")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("content", rootCmd.PersistentFlags().Lookup("content"))
	viper.BindPFlag("status", rootCmd.PersistentFlags().Lookup("status"))
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("content", DefaultContent)
	viper.SetDefault("status", DefaultStatus)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
