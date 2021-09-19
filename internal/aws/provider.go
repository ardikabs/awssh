package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Provider struct {
	Client ec2iface.EC2API
}

func NewProvider(client ec2iface.EC2API) *Provider {
	provider := &Provider{
		Client: client,
	}

	return provider
}

func (p Provider) GetInstanceWithID(instanceID string) ([]*Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	out, err := p.Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	if len(out.Reservations) == 0 {
		return nil, fmt.Errorf("awssh: no instance found")
	}

	instances := p.convert(out.Reservations)
	return instances, nil
}

func (p Provider) GetInstanceWithTag(tags string) ([]*Instance, error) {
	filters, err := PrepareEC2Filters(tags)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	out, err := p.Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	if len(out.Reservations) == 0 {
		return nil, fmt.Errorf("awssh: no instance found")
	}

	instances := p.convert(out.Reservations)
	return instances, nil
}

func (p Provider) convert(ec2Reservations []*ec2.Reservation) []*Instance {
	out := make([]*Instance, 0)

	for i := range ec2Reservations {
		for _, inst := range ec2Reservations[i].Instances {
			out = append(out, NewInstance(inst))
		}
	}
	return out
}
