package aws_test

import (
	. "awssh/internal/aws"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockEC2 struct {
	ec2iface.EC2API

	expectedOutput *ec2.DescribeInstancesOutput
}

func (m *mockEC2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return m.expectedOutput, nil
}

func TestGetInstanceWithID(t *testing.T) {
	expectedOutput := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-12345678abcd"),
						PrivateIpAddress: aws.String("192.168.1.100"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-southeast-1a"),
						},
					},
				},
			},
		},
	}
	provider := NewProvider(&mockEC2{expectedOutput: expectedOutput})

	instance, err := provider.GetInstanceWithID("i-12345678abcd")
	assert.Nil(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, *expectedOutput.Reservations[0].Instances[0].InstanceId, instance[0].InstanceID)
}

func TestGetInstanceWithTag(t *testing.T) {

	t.Run("instance having tag name", func(t *testing.T) {
		expectedOutput := &ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId:       aws.String("i-12345678abcd"),
							PrivateIpAddress: aws.String("192.168.1.100"),
							Placement: &ec2.Placement{
								AvailabilityZone: aws.String("ap-southeast-1a"),
							},
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("lalala"),
								},
							},
						},
					},
				},
			},
		}
		provider := NewProvider(&mockEC2{expectedOutput: expectedOutput})
		instance, err := provider.GetInstanceWithTag("Name=lalala")
		assert.Nil(t, err)
		assert.NotNil(t, instance)
		assert.Equal(t, *expectedOutput.Reservations[0].Instances[0].Tags[0].Value, instance[0].Name)
		assert.Equal(t, "", instance[0].PublicIP)
	})

	t.Run("instance without tag name", func(t *testing.T) {
		expectedOutput := &ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId:       aws.String("i-12345678abcd"),
							PrivateIpAddress: aws.String("192.168.1.100"),
							PublicIpAddress:  aws.String("36.86.63.182"),
							Placement: &ec2.Placement{
								AvailabilityZone: aws.String("ap-southeast-1a"),
							},
							Tags: []*ec2.Tag{},
						},
					},
				},
			},
		}
		provider := NewProvider(&mockEC2{expectedOutput: expectedOutput})
		instance, err := provider.GetInstanceWithTag("Name=lalala")
		assert.Nil(t, err)
		assert.NotNil(t, instance)
		assert.Equal(t, fmt.Sprintf("ec2:noname:%s", *expectedOutput.Reservations[0].Instances[0].InstanceId), instance[0].Name)
		assert.Equal(t, *expectedOutput.Reservations[0].Instances[0].PublicIpAddress, instance[0].PublicIP)
	})
}
