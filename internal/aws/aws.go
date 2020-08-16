package aws

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// NewSession is
func NewSession() *session.Session {
	// TODO: aws.NewSession proper docs
	return session.Must(session.NewSession())
}

func getTagValue(key string, instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func getInstance(session *session.Session, input *ec2.DescribeInstancesInput) (ec2Instances []*EC2Instance, err error) {
	svc := ec2.New(session)
	result, err := svc.DescribeInstances(input)

	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %v", err)
	}

	if len(result.Reservations) == 0 {
		return nil, errors.New("no instance is found")
	}

	reservations := result.Reservations

	for i := range reservations {
		for _, instance := range reservations[i].Instances {
			ec2 := NewEC2Instance(session, instance)
			ec2Instances = append(ec2Instances, ec2)
		}
	}
	return
}

// GetInstanceWithID is
func GetInstanceWithID(session *session.Session, instanceID string) (ec2Instances []*EC2Instance, err error) {
	// TODO: aws.GetInstanceWithID proper docs
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	ec2Instances, err = getInstance(session, input)
	return
}

// GetInstanceWithTag is
func GetInstanceWithTag(session *session.Session, tags string) (ec2Instances []*EC2Instance, err error) {
	// TODO: aws.GetInstanceWithTag proper docs

	input := &ec2.DescribeInstancesInput{
		Filters: prepareFilters(tags),
	}

	ec2Instances, err = getInstance(session, input)
	return
}

func prepareFilters(rawTags string) (filters []*ec2.Filter) {
	awsTags := make(map[string][]*string)

	splitTags := strings.Split(rawTags, ",")

	for _, tags := range splitTags {
		part := strings.Split(tags, "=")
		key := part[0]
		value := aws.String(part[1])

		awsTags[key] = append(awsTags[key], value)
	}

	filters = append(filters, &ec2.Filter{
		Name: aws.String("instance-state-name"),
		Values: []*string{
			aws.String("running"),
		},
	})

	for k, v := range awsTags {
		f := &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", k)),
			Values: v,
		}
		filters = append(filters, f)
	}

	return
}
