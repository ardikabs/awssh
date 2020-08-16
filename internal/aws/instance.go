package aws

import (
	"fmt"
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

// NewEC2Instance is TODO:
func NewEC2Instance(session *session.Session, instance *ec2.Instance) *EC2Instance {
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
	appLogger := config.LoadLogger()

	svc := ec2instanceconnect.New(e.session)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:       aws.String(e.InstanceID),
		SSHPublicKey:     aws.String(publicKey),
		InstanceOSUser:   aws.String(appConfig.SSHUsername),
		AvailabilityZone: aws.String(e.AvailabilityZone),
	}

	appLogger.Debugf("Sending SSH Public Key to the AWS via API", e.InstanceID)

	_, err = svc.SendSSHPublicKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ec2instanceconnect.ErrCodeAuthException:
				return fmt.Errorf("%s: %v", ec2instanceconnect.ErrCodeAuthException, aerr.Error())
			case ec2instanceconnect.ErrCodeInvalidArgsException:
				return fmt.Errorf("%s: %v", ec2instanceconnect.ErrCodeInvalidArgsException, aerr.Error())
			case ec2instanceconnect.ErrCodeServiceException:
				return fmt.Errorf("%s: %v", ec2instanceconnect.ErrCodeServiceException, aerr.Error())
			case ec2instanceconnect.ErrCodeThrottlingException:
				return fmt.Errorf("%s: %v", ec2instanceconnect.ErrCodeThrottlingException, aerr.Error())
			case ec2instanceconnect.ErrCodeEC2InstanceNotFoundException:
				return fmt.Errorf("%s: %v", ec2instanceconnect.ErrCodeEC2InstanceNotFoundException, aerr.Error())
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return fmt.Errorf(err.Error())
		}
	}

	return nil

}

// Connect is a TODO:
func (e *EC2Instance) Connect(usePublicIP bool) (err error) {
	var ipAddr string

	appConfig := config.Get()
	appLogger := config.LoadLogger()
	sshSession, err := ssh.NewSession(e.InstanceID)

	if err != nil {
		return
	}

	err = e.sendSSHPublicKey(sshSession.PublicKey)
	if err != nil {
		return
	}

	ipAddr = e.PrivateIP

	if usePublicIP {
		if e.PublicIP == "" {
			return fmt.Errorf("Could not find public IP for instance %s", e.Name)
		}

		appLogger.Debugf("Use public IP to connect to the EC2 instance target: %s", e.PublicIP)
		ipAddr = e.PublicIP
	}

	appLogger.Debugf("Establish an SSH connection to the EC2 instance target (%s)", e.InstanceID)

	sshArgs := []string{
		"-l",
		appConfig.SSHUsername,
		"-p",
		appConfig.SSHPort,
		ipAddr,
	}

	sshOpts := strings.Split(appConfig.SSHOpts, " ")
	sshArgs = append(sshArgs, sshOpts...)

	fmt.Printf("Running command: ssh %s\n", strings.Join(sshArgs[:], " "))
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
