package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config is a
type Config struct {
	Debug       bool   `envconfig:"debug" default:"0"`
	Tag         string `envconfig:"tag" default:"Name=*"`
	SSHUsername string `envconfig:"ssh_username" default:"ec2-user"`
	SSHPort     string `envconfig:"ssh_port" default:"22"`
	SSHOpts     string `envconfig:"ssh_opts" default:"-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"`
}

var appConfig Config

// Load is used to load the application configuration
func Load() {
	if err := envconfig.Process("awssh", &appConfig); err != nil {
		log.Fatal(err.Error())
	}
}

// Get is used to gather the application configuration
func Get() *Config {
	return &appConfig
}
