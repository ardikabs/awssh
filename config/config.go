package config

import (
	"log"

	"github.com/joeshaw/envdecode"
	flag "github.com/spf13/pflag"
)

// Config represent the application configuration
type config struct {
	Debug       bool   `env:"AWSSH_DEBUG,default=0"`
	Tags        string `env:"AWSSH_TAGS,default=Name=*"`
	SSHUsername string `env:"AWSSH_SSH_USERNAME,default=ec2-user"`
	SSHPort     string `env:"AWSSH_SSH_PORT,default=22"`
	SSHOpts     string `env:"AWSSH_SSH_OPTS,default=-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5"`
	UsePublicIP bool   `env:"AWSSH_USE_PUBLIC_IP,default=0"`
	Region      string `env:"AWS_DEFAULT_REGION"`
}

var appConfig config

// Load used to load the application configuration
func Load() {
	if err := envdecode.Decode(&appConfig); err != nil {
		log.Fatal("Can't load config: ", err)
	}
}

// AddEC2AccessFlags to populate flags used for accessing EC2
func AddEC2AccessFlags(flagSet *flag.FlagSet) {
	flagSet.BoolVarP(&appConfig.Debug, "debug", "d", appConfig.Debug, "Enabled debug mode")
	flagSet.StringVar(&appConfig.Region, "region", appConfig.Region, "Default AWS region to be used. Either set AWS_REGION or AWS_DEFAULT_REGION")
	flagSet.StringVarP(&appConfig.Tags, "tags", "t", appConfig.Tags, "A comma-separated key-value pairs of EC2 tags. Ex: 'Name=ec2,Environment=staging'")
	flagSet.StringVarP(&appConfig.SSHUsername, "ssh-username", "u", appConfig.SSHUsername, "EC2 SSH username")
	flagSet.StringVarP(&appConfig.SSHPort, "ssh-port", "p", appConfig.SSHPort, "An EC2 instance ssh port")
	flagSet.StringVarP(&appConfig.SSHOpts, "ssh-opts", "o", appConfig.SSHOpts, "An additional ssh options")
	flagSet.BoolVarP(&appConfig.UsePublicIP, "use-public-ip", "", appConfig.UsePublicIP, "Use public IP to access the EC2 instance")
}

// GetDebugMode get the debug mode flag
func GetDebugMode() bool {
	return appConfig.Debug
}

// GetRegion get AWS region
func GetRegion() string {
	return appConfig.Region
}

// GetEC2Tags get EC2 tags
func GetEC2Tags() string {
	return appConfig.Tags
}

// GetSSHUsername get SSH username
func GetSSHUsername() string {
	return appConfig.SSHUsername
}

// GetSSHPort get SSH port
func GetSSHPort() string {
	return appConfig.SSHPort
}

// GetSSHOpts get SSH optional argument
func GetSSHOpts() string {
	return appConfig.SSHOpts
}

// GetUsePublicIP get the flag to access EC2 for using public ip or not
func GetUsePublicIP() bool {
	return appConfig.UsePublicIP
}
