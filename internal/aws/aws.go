package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"awssh/internal/logging"
)

// NewSession creates a new AWS session from region input or region environment variables (ex: AWS_DEFAULT_REGION, AWS_REGION)
// all the credentials loaded in a common way of AWS credentials such as,
// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables
// or loaded from AWS shared-credentials located in ~/.aws/credentials
// in particularly when you need to use AWS_PROFILE located in ~/.aws/config
// you need to set AWS_SDK_LOAD_CONFIG=1
func NewSession(region string) (session *aws_session.Session) {
	session = aws_session.Must(aws_session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	logging.Logger().Debugf("Region: %s", *session.Config.Region)

	return
}

// GetInstanceWithID used to get instance with InstanceID input
func GetInstanceWithID(session *aws_session.Session, instanceID string) (ec2Instances []*EC2Instance, err error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	logging.Logger().Debugf("Filter EC2 instances with InstanceID: %s", instanceID)
	ec2Instances, err = getInstance(session, input)

	return
}

// GetInstanceWithTag used to get instance with key-value pair tags input (ex: Environment=production,ProductDomain=VirtualProduct)
func GetInstanceWithTag(session *aws_session.Session, tags string) (ec2Instances []*EC2Instance, err error) {
	filters, err := prepareFilters(tags)

	if err != nil {
		return
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	logging.Logger().Debugf("Filter EC2 instances with tags: %s", tags)
	ec2Instances, err = getInstance(session, input)

	return
}

// prepareFilters used to form a proper AWS filter tags format
// from a raw tags input format
func prepareFilters(rawTags string) (filters []*ec2.Filter, err error) {
	awsTags := make(map[string][]*string)

	splitTags := strings.Split(rawTags, ",")

	for _, tags := range splitTags {
		part := strings.Split(tags, "=")

		if len(part) != 2 {
			return nil, fmt.Errorf("Wrong tag format, it should be follow 'Key=Value' format: '%s'", tags)
		}

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

	logging.Logger().Debugf("Use the following filters to filter EC2 instances: %v", filters)

	return
}

// getInstance handle the underlaying to gather the EC2 instances
// following with the DescribeInstances method
func getInstance(session *aws_session.Session, input *ec2.DescribeInstancesInput) (ec2Instances []*EC2Instance, err error) {
	svc := ec2.New(session)
	result, err := svc.DescribeInstances(input)

	if err != nil {
		return nil, fmt.Errorf("Failed to get instance: (%v)", err)
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("No instance is found on region %s", *session.Config.Region)
	}

	reservations := result.Reservations

	for i := range reservations {
		for _, instance := range reservations[i].Instances {
			ec2 := NewEC2Instance(session, instance)
			ec2Instances = append(ec2Instances, ec2)
		}
	}

	logging.Logger().Debugf("Found %d EC2 instances on region %s", len(ec2Instances), *session.Config.Region)

	return
}

func getTagValue(key string, instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
