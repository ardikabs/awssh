package cli

import (
	"awssh/config"
	"awssh/internal/aws"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	usePublicIP bool
)

func init() {
	appConfig := config.Get()

	rootCommand.Flags().BoolVarP(&appConfig.Debug, "debug", "d", false, "Enabled debug mode")
	rootCommand.Flags().StringVarP(&appConfig.Tag, "tag", "t", "Name=*", "EC2 tag keypair")
	rootCommand.Flags().StringVarP(&appConfig.SSHUsername, "ssh-username", "u", "ec2-user", "EC2 SSH username")
	rootCommand.Flags().StringVarP(&appConfig.SSHPort, "ssh-port", "p", "22", "An EC2 instance ssh port")
	rootCommand.Flags().StringVarP(&appConfig.SSHOpts, "ssh-opts", "o", "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null", "An ssh optional arguments")
	rootCommand.Flags().BoolVarP(&usePublicIP, "use-public-ip", "", false, "Use public IP to access the EC2 instance")
}

var rootCommand = &cobra.Command{
	Use:   "awssh",
	Short: "awssh is a simple CLI to ssh'ing EC2",
	Long:  "awssh is a simple CLI providing an ssh access to EC2 utilizing ec2-instance-connect",
	Run:   runSSHAccess,
}

func runSSHAccess(cmd *cobra.Command, args []string) {
	appConfig := config.Get()

	// get AWS session in common way, env variables and shared-credential file
	session := aws.NewSession()

	instances, err := aws.GetInstanceWithTag(session, appConfig.Tag)

	if err != nil {
		log.Fatal(err)
	}

	instance, err := promptUI(instances)

	if err != nil {
		// TODO: Fatal log on CLI
		log.Fatal(err)
	}

	instance.Connect(usePublicIP)
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
