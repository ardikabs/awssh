package main

import (
	"awssh/config"
	"awssh/internal/cli"
	"os"
)

func main() {

	config.Load()

	rootCommand := cli.MakeRoot()

	rootCommand.AddCommand(cli.MakeVersion())

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

// TODO: unittest
