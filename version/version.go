package version

import "fmt"

// Version represent awssh current version
var Version = getLocalVersion()

func getLocalVersion() string {
	return fmt.Sprintf("v0.1.0")
}
