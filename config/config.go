package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

// Config represent the application configuration
type Config struct {
	Debug       bool   `env:"AWSSH_DEBUG,default=0"`
	Tags        string `env:"AWSS_TAGS,default=Name=*"`
	SSHUsername string `env:"AWSSH_SSH_USERNAME,default=ec2-user"`
	SSHPort     string `env:"AWSSH_SSH_PORT,default=22"`
	SSHOpts     string `env:"AWSSH_SSH_OPTS,default=-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5"`
	Region      string `env:"AWS_DEFAULT_REGION"`
}

var appConfig Config

// Load used to load the application configuration
func Load() {

	if err := envdecode.Decode(&appConfig); err != nil {
		log.Fatal("Can't load config: ", err)
	}
}

// Get used to get the application configuration
func Get() *Config {
	return &appConfig
}
