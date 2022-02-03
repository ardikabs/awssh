package aws

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"golang.org/x/crypto/ssh/agent"

	"awssh/config"
	"awssh/internal/logging"
	"awssh/internal/ssh"
)

// EC2Instance represent all the necessary EC2 components
type Instance struct {
	Name             string
	InstanceID       string
	PrivateIP        string
	PublicIP         string
	AvailabilityZone string
}

type ShellCommandFunc func(name string, args ...string) *exec.Cmd

// NewEC2Instance creates a new EC2Instance from aws ec2 instance source
func NewInstance(instance *ec2.Instance) *Instance {
	ec2InstanceName := GetTagValue("Name", instance)

	if ec2InstanceName == "" {
		ec2InstanceName = fmt.Sprintf("ec2:noname:%s", *instance.InstanceId)
	}

	var publicIPAddr string

	if instance.PublicIpAddress == nil {
		publicIPAddr = ""
	} else {
		publicIPAddr = *instance.PublicIpAddress
	}

	return &Instance{
		Name:             ec2InstanceName,
		InstanceID:       *instance.InstanceId,
		PrivateIP:        *instance.PrivateIpAddress,
		PublicIP:         publicIPAddr,
		AvailabilityZone: *instance.Placement.AvailabilityZone,
	}
}

// sendSSHPublicKey is an extend method to do ec2-instance-connect task
// for sending SSH Public Key to the AWS API Server
func (e *Instance) sendSSHPublicKey(client ec2instanceconnectiface.EC2InstanceConnectAPI, publicKey string) (err error) {
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:       aws.String(e.InstanceID),
		SSHPublicKey:     aws.String(publicKey),
		InstanceOSUser:   aws.String(config.GetSSHUsername()),
		AvailabilityZone: aws.String(e.AvailabilityZone),
	}

	logging.Logger().Debugf("Sending SSH Public Key for EC2 instance '%s' (%s)", e.Name, e.InstanceID)

	if _, err := client.SendSSHPublicKey(input); err != nil {
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
func (e *Instance) Connect(sshAgent agent.ExtendedAgent, client ec2instanceconnectiface.EC2InstanceConnectAPI, cmdFn ShellCommandFunc, usePublicIP bool) (err error) {
	logging.Logger().Debugf("awssh: select EC2 instance '%s' (%s)", e.Name, e.InstanceID)

	var ipAddr string

	sshSession, err := ssh.NewSession(sshAgent, e.InstanceID)
	if err != nil {
		return
	}

	if err := e.sendSSHPublicKey(client, sshSession.PublicKey); err != nil {
		return err
	}

	ipAddr = e.PrivateIP

	if usePublicIP {
		if e.PublicIP == "" {
			return fmt.Errorf("awssh: could not find public IP for EC2 instance target '%s' (%s)", e.Name, e.InstanceID)
		}

		logging.Logger().Debugf("awssh: use public IP to connect to the EC2 instance target '%s' (%s): %s", e.Name, e.InstanceID, e.PublicIP)
		ipAddr = e.PublicIP
	}

	logging.Logger().Debugf("awssh: establish an SSH connection to the EC2 instance target '%s' (%s)", e.Name, e.InstanceID)

	sshArgs := []string{
		"-l",
		config.GetSSHUsername(),
		"-p",
		config.GetSSHPort(),
		ipAddr,
	}

	sshOpts := strings.Split(config.GetSSHOpts(), " ")
	sshArgs = append(sshArgs, sshOpts...)

	logging.Logger().Infof("awssh: running command: ssh %s\n", strings.Join(sshArgs[:], " "))
	return cmdFn("ssh", sshArgs...).Run()
}
