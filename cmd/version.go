package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// VersionGitCommit The git commit this build is from
var VersionGitCommit string

// VersionGitTag The git tag this build is from
var VersionGitTag string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Convey",
	Run: func(cmd *cobra.Command, args []string) {
		if VersionGitCommit != "" {
			fmt.Printf("Convey %s -- %s\n\n", VersionGitTag, VersionGitCommit[:7])
			fmt.Printf("Commit: %s\n", VersionGitCommit)
		}
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Go version: %s\n", runtime.Version())
	},
}
