package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"awssh/internal/logging"
)

// NewSession creates a new AWS session from region input or region environment variables (ex: AWS_DEFAULT_REGION, AWS_REGION)
// all the credentials loaded in a common way of AWS credentials such as,
// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables
// or loaded from AWS shared-credentials located in ~/.aws/credentials
// in particularly when you need to use AWS_PROFILE located in ~/.aws/config
// you need to set AWS_SDK_LOAD_CONFIG=1
//
// Sidenote
// session.Must(): the only way the session is failed if the shared config is malformed
// ref: https://github.com/aws/aws-sdk-go/issues/928
func NewSession(region string) *session.Session {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	logging.Logger().Debugf("Region: %s", *session.Config.Region)

	return session
}

func PrepareEC2Filters(tags string) ([]*ec2.Filter, error) {
	awsTags := make(map[string][]*string)

	splitTags := strings.Split(tags, ",")

	for _, tags := range splitTags {
		part := strings.Split(tags, "=")

		if len(part) != 2 {
			return nil, fmt.Errorf("awssh: bad input, filters must be using 'Key=Value' format: '%s'", tags)
		}

		key := part[0]
		value := aws.String(part[1])
		awsTags[key] = append(awsTags[key], value)
	}

	filters := make([]*ec2.Filter, 0)
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

	return filters, nil
}

func GetTagValue(key string, instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}

	return ""
}
