package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(versionCmd)
}

// Version represent an awssh version
var (
	Version   string
	GitCommit string
)

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Print the version number of awssh",
	Long:         `All software has versions. This is awssh's`,
	Example:      `  awssh version`,
	SilenceUsage: false,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Awssh an extended ec2-instance-connect command tool")
		printVersion()
	},
}

func printVersion() {
	if len(Version) == 0 {
		fmt.Println("Version: dev")
	} else {
		fmt.Println("Version:", Version)
	}
	fmt.Println("Git Commit:", GitCommit)
}
