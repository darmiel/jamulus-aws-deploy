package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// headers
const (
	JamulusDefHeader   = "Jamulus"
	JamulusStateHeader = "Jamulus-State"
	JamulusOwnerHeader = "Jamulus-Owner"
)

// values
const (
	JamulusStateCreated = "created"
	JamulusStateDone    = "done"
)

var (
	JamulusDefTag = ec2.Tag{
		Key:   aws.String(JamulusDefHeader),
		Value: aws.String("Yes"),
	}
	JamulusStatusCreatedTag = ec2.Tag{
		Key:   aws.String(JamulusStateHeader),
		Value: aws.String(JamulusStateCreated),
	}
)

func Is(instance *ec2.Instance, state string) bool {
	for _, tag := range instance.Tags {
		if tag.Key != nil && *tag.Key == JamulusStateHeader {
			if tag.Value != nil && *tag.Value == state {
				return true
			}
		}
	}
	return false
}
