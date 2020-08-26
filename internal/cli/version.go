package cli

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

// Version represent an awssh version
var (
	Version   string
	GitCommit string
)

const awsshFigletStr = `                        _
                       | |
  __ ___      _____ ___| |__
 / _` + "`" + ` \ \ /\ / / __/ __| '_ \
| (_| |\ V  V /\__ \__ \ | | |
 \__,_| \_/\_/ |___/___/_| |_|

an extended ec2-instance-connect command tool

`

// MakeVersion used to create version subcommand
func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version number of awssh",
		Long:         `All software has versions. This is awssh's`,
		Example:      `  awssh version`,
		SilenceUsage: false,
	}

	command.Run = func(cmd *cobra.Command, args []string) {
		printVersion()
	}

	return command
}

func printVersion() {
	printASCIIArt()
	if len(Version) == 0 {
		fmt.Println("Version: dev")
	} else {
		fmt.Println("Version:", Version)
	}
	fmt.Println("Git Commit:", GitCommit)
}

func printASCIIArt() {
	awsshLogo := aec.CyanF.Apply(awsshFigletStr)
	fmt.Print(awsshLogo)
}
