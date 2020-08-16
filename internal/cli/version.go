package cli

import (
	"fmt"

	"awssh/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of awssh",
	Long:  `All software has versions. This is awssh's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("awssh %s\n", version.Version)
	},
}
