package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var natsUrl string
var natsClusterId string
var forceWrite bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.PersistentFlags().StringVar(&natsUrl, "nats-url", "", "NATS server url")
	configureCmd.PersistentFlags().StringVar(&natsClusterId, "nats-cluster", "", "NATS cluster id")
	configureCmd.PersistentFlags().BoolVar(&forceWrite, "force", false, "Overwrite current configuration")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Convey",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("NatsURL", natsUrl)
		viper.SetDefault("NatsClusterID", natsClusterId)
		if forceWrite {
			err := viper.WriteConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			err := viper.SafeWriteConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
