package cmd

import (
	"fmt"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The NATS URL passed in from command-line
var natsURL string

// The NATS cluster ID passed in from command-line
var natsClusterID string

// Whether short channel names should be used instead of the standard uuid format
var useShortName bool

// Whether the current config file should be overwritten
var forceWrite bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.PersistentFlags().StringVar(&natsURL, "nats-url", "", "NATS server url")
	configureCmd.PersistentFlags().StringVar(&natsClusterID, "nats-cluster", "", "NATS cluster id")
	configureCmd.PersistentFlags().BoolVar(&useShortName, "short-name", false, "Use short channel names. Channel conflicts may occur.")
	configureCmd.PersistentFlags().BoolVar(&forceWrite, "overwrite", false, "Overwrite current configuration")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Convey",
	Run:   ConfigureCommandFunc,
}

// ConfigureCommandFunc is a handler for the configure command
func ConfigureCommandFunc(cmd *cobra.Command, args []string) {
	viper.Set(configKeyNatsURL, natsURL)
	viper.Set(configKeyNatsClusterID, natsClusterID)
	viper.Set(configKeyUseShortName, useShortName)

	// If a config file is found, read it in.
	configFileExists := false
	if err := viper.ReadInConfig(); err == nil {
		configFileExists = true
	}

	// If config file doesn't exist and it hasn't been set in viper, set it
	if !configFileExists && viper.ConfigFileUsed() == "" {
		home, err := homedir.Dir()
		if err != nil {
			errorExit(err.Error())
		}
		viper.SetConfigFile(path.Join(home, ".convey.yaml"))
	}

	configFilePath := viper.ConfigFileUsed()

	if forceWrite || !configFileExists {
		err := viper.WriteConfigAs(configFilePath)
		if err != nil {
			errorExit(err.Error())
		}
	} else {
		msg := fmt.Sprintf("Config file exists. Use --overwrite to overwrite the config file at %s", configFilePath)
		errorExit(msg)
	}
}
