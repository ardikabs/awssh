package aws_test

import (
	"fmt"
	"testing"

	. "awssh/internal/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	session := NewSession("ap-southeast-1")
	assert.NotNil(t, session)
}

func TestPrepareEC2Filters(t *testing.T) {
	type args struct {
		tags string
	}

	defaultFilters := []*ec2.Filter{
		{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		},
	}

	tests := []struct {
		name     string
		args     args
		expected []*ec2.Filter
		err      error
	}{
		{
			name:     "no argument (empty)",
			args:     args{""},
			expected: nil,
			err:      fmt.Errorf("awssh: bad input, filters must be using 'Key=Value' format: ''"),
		},
		{
			name:     "wrong format argument",
			args:     args{"Environment,production"},
			expected: nil,
			err:      fmt.Errorf("awssh: bad input, filters must be using 'Key=Value' format: 'Environment'"),
		},
		{
			name: "single tag",
			args: args{"Environment=production"},
			expected: append(defaultFilters,
				&ec2.Filter{
					Name: aws.String("tag:Environment"),
					Values: []*string{
						aws.String("production"),
					},
				},
			),
			err: nil,
		},
		{
			name: "multiple tags with comma delimiters",
			args: args{"Environment=production,Service=promotion"},
			expected: append(defaultFilters,
				&ec2.Filter{
					Name: aws.String("tag:Environment"),
					Values: []*string{
						aws.String("production"),
					},
				},
				&ec2.Filter{
					Name: aws.String("tag:Service"),
					Values: []*string{
						aws.String("promotion"),
					},
				},
			),
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareEC2Filters(tt.args.tags)
			if tt.err != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}

func TestGetTagValue(t *testing.T) {
	defaultInstance := &ec2.Instance{
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Environment"),
				Value: aws.String("production"),
			},
		},
	}

	type args struct {
		key      string
		instance *ec2.Instance
	}

	tests := []struct {
		name     string
		args     args
		expected string
		err      error
	}{
		{
			name:     "instance with no specified tags",
			args:     args{"ManagedBy", defaultInstance},
			expected: "",
			err:      nil,
		},
		{
			name:     "instance with specified tags",
			args:     args{"Environment", defaultInstance},
			expected: "production",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := GetTagValue(tt.args.key, tt.args.instance)
			assert.Equal(t, tt.expected, want)
		})
	}
}
