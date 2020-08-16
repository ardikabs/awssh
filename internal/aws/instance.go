package aws

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"awssh/config"
	"awssh/internal/ssh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
)

// EC2Instance is
type EC2Instance struct {
	session *session.Session

	Name             string
	InstanceID       string
	PrivateIP        string
	PublicIP         string
	AvailabilityZone string
}

// NewEC2Instance is
func NewEC2Instance(session *session.Session, instance *ec2.Instance) *EC2Instance {
	// TODO: aws.NewEC2Instance proper docs
	ec2InstanceName := getTagValue("Name", instance)

	var publicIPAddr string

	if instance.PublicIpAddress == nil {
		publicIPAddr = ""
	} else {
		publicIPAddr = *instance.PublicIpAddress
	}

	return &EC2Instance{
		session:          session,
		Name:             ec2InstanceName,
		InstanceID:       *instance.InstanceId,
		PrivateIP:        *instance.PrivateIpAddress,
		PublicIP:         publicIPAddr,
		AvailabilityZone: *instance.Placement.AvailabilityZone,
	}
}

func (e *EC2Instance) sendSSHPublicKey(publicKey string) (err error) {
	appConfig := config.Get()

	svc := ec2instanceconnect.New(e.session)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:       aws.String(e.InstanceID),
		AvailabilityZone: aws.String(e.AvailabilityZone),
		InstanceOSUser:   aws.String(appConfig.SSHUsername),
		SSHPublicKey:     aws.String(publicKey),
	}

	_, err = svc.SendSSHPublicKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ec2instanceconnect.ErrCodeAuthException:
				fmt.Println(ec2instanceconnect.ErrCodeAuthException, aerr.Error())
			case ec2instanceconnect.ErrCodeInvalidArgsException:
				fmt.Println(ec2instanceconnect.ErrCodeInvalidArgsException, aerr.Error())
			case ec2instanceconnect.ErrCodeServiceException:
				fmt.Println(ec2instanceconnect.ErrCodeServiceException, aerr.Error())
			case ec2instanceconnect.ErrCodeThrottlingException:
				fmt.Println(ec2instanceconnect.ErrCodeThrottlingException, aerr.Error())
			case ec2instanceconnect.ErrCodeEC2InstanceNotFoundException:
				fmt.Println(ec2instanceconnect.ErrCodeEC2InstanceNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	return nil

}

// Connect is a
func (e *EC2Instance) Connect(usePublicIP bool) (err error) {
	// TODO: aws.Connect proper docs

	var ipAddr string

	appConfig := config.Get()
	sshSession, err := ssh.NewSession(e.InstanceID)

	if err != nil {
		log.Fatalf("unable to raise an ssh session: %v", err)
	}

	err = e.sendSSHPublicKey(sshSession.PublicKey)
	if err != nil {
		log.Fatalf("unable to send public key to AWS: %v", err)
	}

	ipAddr = e.PrivateIP

	if usePublicIP {
		if e.PublicIP == "" {
			return fmt.Errorf("Could not find public IP for instance %s", e.Name)
		}
		ipAddr = e.PublicIP
	}

	sshArgs := []string{
		"-l",
		appConfig.SSHUsername,
		"-p",
		appConfig.SSHPort,
		ipAddr,
	}

	sshArgs = append(sshArgs, strings.Split(appConfig.SSHOpts, " ")...)

	fmt.Printf("Running command: ssh %s\n", strings.Join(sshArgs[:], " "))

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
