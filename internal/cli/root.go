package cli

import (
	"awssh/config"
	"awssh/internal/aws"
	"awssh/internal/logging"

	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	usePublicIP bool
	region      string
)

func init() {
	appConfig := config.Get()

	rootCommand.Flags().BoolVarP(&appConfig.Debug, "debug", "d", false, "Enabled debug mode")
	rootCommand.Flags().StringVarP(&appConfig.Tags, "tags", "t", "Name=*", "EC2 tags key-value pair")
	rootCommand.Flags().StringVarP(&appConfig.SSHUsername, "ssh-username", "u", "ec2-user", "EC2 SSH username")
	rootCommand.Flags().StringVarP(&appConfig.SSHPort, "ssh-port", "p", "22", "An EC2 instance ssh port")
	rootCommand.Flags().StringVarP(&appConfig.SSHOpts, "ssh-opts", "o", "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=5", "An additional ssh options")
	rootCommand.Flags().BoolVarP(&usePublicIP, "use-public-ip", "", false, "Use public IP to access the EC2 instance")
	rootCommand.Flags().StringVarP(&region, "region", "", os.Getenv("AWS_DEFAULT_REGION"), "Default AWS region to be used. Either set AWS_REGION or AWS_DEFAULT_REGION")
}

var rootCommand = &cobra.Command{
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

  #  public ip to connect to the EC2 instance
  awssh --use-public-ip
`,
	Args: validateInstanceIDArgs,
	Run:  runSSHAccess,
}

func validateInstanceIDArgs(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		match, _ := regexp.MatchString(`^i-[\w]+`, args[0])
		if !match {
			return fmt.Errorf("Invalid instance-id format: '%s'", args[0])
		}
	}
	return
}

func runSSHAccess(cmd *cobra.Command, args []string) {
	var target *aws.EC2Instance

	appConfig := config.Get()

	logging.Init(appConfig.Debug)

	session := aws.NewSession(region)

	if len(args) > 0 {
		instances, err := aws.GetInstanceWithID(session, args[0])
		if err != nil {
			exitWithError(err)
		}

		target = instances[0]
	} else {
		instances, err := aws.GetInstanceWithTag(session, appConfig.Tags)
		if err != nil {
			exitWithError(err)
		}

		target, err = promptUI(instances)
		if err != nil {
			exitWithError(err)
		}
	}

	err := target.Connect(usePublicIP)
	if err != nil {
		exitWithError(err)
	}
}

func promptUI(instances []*aws.EC2Instance) (instance *aws.EC2Instance, err error) {
	searcher := func(i string, index int) bool {
		inst := instances[index]
		name := inst.Name
		input := i
		return strings.Contains(name, input) || strings.Contains(inst.InstanceID, input) || strings.Contains(inst.PrivateIP, input) || strings.Contains(inst.PrivateIP, input)
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
