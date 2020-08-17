package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config represent the application configuration
type Config struct {
	Debug       bool   `envconfig:"debug" default:"0"`
	Tags        string `envconfig:"tags" default:"Name=*"`
	SSHUsername string `envconfig:"ssh_username" default:"ec2-user"`
	SSHPort     string `envconfig:"ssh_port" default:"22"`
	SSHOpts     string `envconfig:"ssh_opts" default:"-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5"`
}

var appConfig Config

// Load used to load the application configuration
func Load() {
	if err := envconfig.Process("awssh", &appConfig); err != nil {
		log.Fatal("Can't load config: ", err)
	}
}

// Get used to get the application configuration
func Get() *Config {
	return &appConfig
}
