package aws

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"

	"awssh/config"
	"awssh/internal/logging"
	"awssh/internal/ssh"
)

// EC2Instance represent all the necessary EC2 components
type EC2Instance struct {
	session *aws_session.Session

	Name             string
	InstanceID       string
	PrivateIP        string
	PublicIP         string
	AvailabilityZone string
}

// NewEC2Instance creates a new EC2Instance from aws ec2 instance source
func NewEC2Instance(session *aws_session.Session, instance *ec2.Instance) *EC2Instance {
	ec2InstanceName := getTagValue("Name", instance)

	if ec2InstanceName == "" {
		ec2InstanceName = fmt.Sprintf("ec2:noname:%s", *instance.InstanceId)
	}

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

// sendSSHPublicKey is an extend method to do ec2-instance-connect task
// for sending SSH Public Key to the AWS API Server
func (e *EC2Instance) sendSSHPublicKey(publicKey string) (err error) {
	appConfig := config.Get()
	logger := logging.Get()

	svc := ec2instanceconnect.New(e.session)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:       aws.String(e.InstanceID),
		SSHPublicKey:     aws.String(publicKey),
		InstanceOSUser:   aws.String(appConfig.SSHUsername),
		AvailabilityZone: aws.String(e.AvailabilityZone),
	}

	logger.Debugf("Sending SSH Public Key for EC2 instance '%s' (%s)", e.Name, e.InstanceID)

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

// Connect used to establish an ssh connection from the EC2Instance
// following with the use of public ip
func (e *EC2Instance) Connect(usePublicIP bool) (err error) {
	var ipAddr string

	appConfig := config.Get()
	logger := logging.Get()

	logger.Debugf("Select EC2 instance '%s' (%s)", e.Name, e.InstanceID)

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
			return fmt.Errorf("Could not find public IP for EC2 instance target '%s' (%s)", e.Name, e.InstanceID)
		}

		logger.Debugf("Use public IP to connect to the EC2 instance target '%s' (%s): %s", e.Name, e.InstanceID, e.PublicIP)
		ipAddr = e.PublicIP
	}

	logger.Debugf("Establish an SSH connection to the EC2 instance target '%s' (%s)", e.Name, e.InstanceID)

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
