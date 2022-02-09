package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"github.com/lukasmoellerch/mensa-cli/internal/eth"
	"github.com/lukasmoellerch/mensa-cli/internal/uzh"
	"github.com/mailru/easyjson"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfg config

var cfgFile string
var langFlag string
var storageDirectory string
var cpuProfile string

var providers []base.Provider

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mensa-cli",
	Short: "A CLI tool which fetches the list of meals for a given date",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&storageDirectory, "storage-directory", "", "The directory where the data is stored (default is $HOME/.mensa-cli)")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is storage-directory/.mensa-cli/config.yaml)")
	RootCmd.PersistentFlags().StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile to file")
	RootCmd.PersistentFlags().StringVar(&langFlag, "lang", "en", "language to use")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if storageDirectory == "" {
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		storageDirectory = path.Join(home, ".mensa-cli")
	}

	// Create Folder if it doesn't exist
	if _, err := os.Stat(storageDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(storageDirectory, 0755)
		cobra.CheckErr(err)
	}

	if cfgFile == "" {
		cfgFile = path.Join(storageDirectory, "config.json")
	}

	// Read in config file manually using io
	f, err := os.Open(cfgFile)
	if err != nil {
		fmt.Println("Config file not found, using default values")
	} else {
		easyjson.UnmarshalFromReader(f, &cfg)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	uzhProvider := uzh.Provider{}
	ethProvider := eth.Provider{}

	providers = []base.Provider{&uzhProvider, &ethProvider}
}

func writeConfig() error {
	f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	} else {
		_, err := easyjson.MarshalToWriter(&cfg, f)
		if err != nil {
			return err
		}
	}
	return nil
}
