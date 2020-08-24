package cli

import (
	"awssh/config"
	"fmt"
	"os"
)

// Execute is an entrypoint for the command line interface
func Execute() {
	config.Load()

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

// exitWithError will terminate execution with an error result
// It prints the error to stderr and exits with a non-zero exit code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
