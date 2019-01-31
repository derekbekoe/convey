package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// TODO-DEREK Add commit hash to version - https://blog.alexellis.io/inject-build-time-vars-golang/

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Convey",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Convey v0.0.1 -- HEAD")
	},
}
