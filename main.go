package main

import (
	"os"

	"awssh/cmd"
	"awssh/config"
)

func main() {

	config.Load()

	rootCmd := cmd.MakeRoot()
	versionCmd := cmd.MakeVersion()

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
