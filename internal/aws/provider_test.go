package aws_test

import (
	. "awssh/internal/aws"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockEC2 struct {
	ec2iface.EC2API
}

func (*mockEC2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
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
	}, nil
}

func TestGetInstanceWithID(t *testing.T) {
	provider := NewProvider(&mockEC2{})

	instance, err := provider.GetInstanceWithID("i-12345678abcd")
	assert.Equal(t, "i-12345678abcd", instance[0].InstanceID)
	assert.Nil(t, err)
	assert.NotNil(t, instance)
}

func TestGetInstanceWithTag(t *testing.T) {
	provider := NewProvider(&mockEC2{})

	instance, err := provider.GetInstanceWithTag("Name=lalala")
	assert.Nil(t, err)
	assert.NotNil(t, instance)
	assert.Equal(t, "lalala", instance[0].Name)
}
