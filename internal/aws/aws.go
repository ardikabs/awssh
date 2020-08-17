package aws

import (
	"awssh/config"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// NewSession creates a new AWS session from region input or region environment variables (ex: AWS_DEFAULT_REGION, AWS_REGION)
// all the credentials loaded in a common way of AWS credentials such as,
// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables
// or loaded from AWS shared-credentials located in ~/.aws/credentials
// in particularly when you need to use AWS_PROFILE located in ~/.aws/config
// you need to set AWS_SDK_LOAD_CONFIG=1
func NewSession(region string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
}

// GetInstanceWithID used to get instance with InstanceID input
func GetInstanceWithID(session *session.Session, instanceID string) (ec2Instances []*EC2Instance, err error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	ec2Instances, err = getInstance(session, input)
	return
}

// GetInstanceWithTag used to get instance with key-value pair tags input (ex: Environment=production,ProductDomain=VirtualProduct)
func GetInstanceWithTag(session *session.Session, tags string) (ec2Instances []*EC2Instance, err error) {

	input := &ec2.DescribeInstancesInput{
		Filters: prepareFilters(tags),
	}

	ec2Instances, err = getInstance(session, input)
	return
}

// prepareFilters used to form a proper AWS filter tags format
// from a raw tags input format
func prepareFilters(rawTags string) (filters []*ec2.Filter) {
	appLogger := config.LoadLogger()

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

	appLogger.Debugf("Use the following filters to filter EC2 instances: %v", filters)

	return
}

// getInstance handle the underlaying to gather the EC2 instances
// following with the DescribeInstances method
func getInstance(session *session.Session, input *ec2.DescribeInstancesInput) (ec2Instances []*EC2Instance, err error) {
	svc := ec2.New(session)
	result, err := svc.DescribeInstances(input)

	if err != nil {
		return nil, fmt.Errorf("Failed to get instance: (%v)", err)
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("No instance is found")
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

func getTagValue(key string, instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
