package aws

import (
	"awssh/config"
	"awssh/internal/logging"
	"os"
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh/agent"
)

type mockEC2InstanceConnectAPI struct {
	ec2instanceconnectiface.EC2InstanceConnectAPI

	expectedInput *ec2instanceconnect.SendSSHPublicKeyInput
}

func (m mockEC2InstanceConnectAPI) SendSSHPublicKey(input *ec2instanceconnect.SendSSHPublicKeyInput) (*ec2instanceconnect.SendSSHPublicKeyOutput, error) {

	if *(m.expectedInput.InstanceId) != *(input.InstanceId) {
		return nil, awserr.New(ec2instanceconnect.ErrCodeInvalidArgsException, "mismatch instance-id", nil)
	}

	return nil, nil
}

type mockSSHAgent struct {
	agent.ExtendedAgent
}

func (m mockSSHAgent) List() ([]*agent.Key, error) {
	return []*agent.Key{}, nil
}

func (m mockSSHAgent) Add(key agent.AddedKey) error {
	return nil
}

func TestMain(m *testing.M) {
	config.Load()
	logging.NewLogger(false)
	code := m.Run()
	os.Exit(code)
}

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	os.Exit(0)
}

func fakeShellCommand() func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{
			"-test.run=TestShellProcessSuccess",
			"--",
			name,
		}

		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_PROCESS=1"}
		return cmd
	}
}

func TestConnect(t *testing.T) {
	defaultInstance := &ec2.Instance{
		InstanceId:       aws.String("i-1234567890"),
		PrivateIpAddress: aws.String("10.10.5.100"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Environment"),
				Value: aws.String("production"),
			},
		},
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String("ap-southeast-1a"),
		},
	}

	instance := NewInstance(defaultInstance)
	mockEC2InstanceConnectAPI := mockEC2InstanceConnectAPI{
		expectedInput: &ec2instanceconnect.SendSSHPublicKeyInput{
			InstanceId: aws.String("i-1234567890"),
		},
	}

	shellCommand := fakeShellCommand()

	mockSSHAgent := mockSSHAgent{}
	err := instance.Connect(mockSSHAgent, mockEC2InstanceConnectAPI, shellCommand, false)
	assert.Nil(t, err)
}
