package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"awssh/config"
	"awssh/internal/aws"
	"awssh/internal/logging"
)

// MakeRoot used to create a root command functionality
func MakeRoot() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "awssh",
		Short: "awssh is a simple CLI to ssh'ing EC2",
		Long:  "awssh is a simple CLI providing an ssh access to EC2 utilizing ec2-instance-connect",
		Example: `
	  # List all of the EC2 instances given by the credentials
	  awssh --region=ap-southeast-1

	  # Select EC2 instance with instance-id
	  awssh i-0387e016c47c6170c

	  # Select EC2 instance given with selected tags
	  awssh --tags "Environment=production,Project=jenkins,Owner=SRE"

	  # Use an additional ssh options
	  awssh --tags "Environment=staging,ProductDomain=bastion" --ssh-username=centos --ssh-port=2222 --ssh-opts="-o ServerAliveInterval=60s"

	  # Use public ip to connect to the EC2 instance
	  awssh --use-public-ip
	`,
	}

	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Run = runSSHAccess

	config.AddEC2AccessFlags(cmd.Flags())
	return cmd
}

func validateInstanceIDArgs(args []string) (err error) {
	if len(args) > 0 {
		match, _ := regexp.MatchString(`^i-[\w]+`, args[0])
		if !match {
			return fmt.Errorf("Invalid instance-id format: '%s'", args[0])
		}
	}
	return
}

func runSSHAccess(cmd *cobra.Command, args []string) {
	logging.NewLogger(config.GetDebugMode())

	err := validateInstanceIDArgs(args)

	if err != nil {
		logging.ExitWithError(err)
	}

	var target *aws.EC2Instance

	session := aws.NewSession(config.GetRegion())

	if len(args) > 0 {
		instances, err := aws.GetInstanceWithID(session, args[0])
		if err != nil {
			logging.ExitWithError(err)
		}

		target = instances[0]
	} else {
		instances, err := aws.GetInstanceWithTag(session, config.GetEC2Tags())
		if err != nil {
			logging.ExitWithError(err)
		}

		target, err = promptUI(instances)
		if err != nil {
			logging.ExitWithError(err)
		}
	}

	err = target.Connect(config.GetUsePublicIP())
	if err != nil {
		logging.ExitWithError(err)
	}
}

func promptUI(instances []*aws.EC2Instance) (instance *aws.EC2Instance, err error) {
	searcher := func(i string, index int) bool {
		inst := instances[index]
		name := inst.Name
		input := i
		return strings.Contains(name, input) || strings.Contains(inst.InstanceID, input) || strings.Contains(inst.PrivateIP, input) || strings.Contains(inst.PublicIP, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}`,
		Active:   `{{ "»" | magenta }} {{ .Name | yellow }} {{ .InstanceID | green }} ({{ .PrivateIP | red }}{{if ne .PublicIP "" }} {{"/"}} {{ .PublicIP | red }}{{ end }})`,
		Inactive: `  {{ .Name }} {{ .InstanceID | cyan }} ({{ .PrivateIP }}{{if ne .PublicIP "" }} {{"/"}} {{ .PublicIP }}{{ end }})`,
		Selected: `{{ .Name | green }} {{ .InstanceID | red }}`,
	}

	prompt := &promptui.Select{
		Label:     "Select an instance:",
		Items:     instances,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return instances[i], nil
}
